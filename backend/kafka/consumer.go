package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/websockets"
	// "github.com/gocql/gocql"
)

var hub *websockets.Hub

// StartConsumer continuously logs incoming messages
func StartConsumer(ctx context.Context) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}
	defer client.Close()

	partitionConsumer, err := client.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume Kafka topic: %v", err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Kafka Consumer started...")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Kafka message: %s\n", string(msg.Value))
			var message models.Message
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
			}

		case <-ctx.Done():
			fmt.Println("Consumer shutting down...")
			return
		}
	}
}

func WebsocketConsumer(ctx context.Context, websocketHub *websockets.Hub) {
	hub = websocketHub
	keepRunning := true
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}


	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, "websocket-group", config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for { 
			err := client.Consume(ctx, []string{topic}, &consumer)
			if err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Printf("Error from consumer: %v", err)
			}

			// If context is canceled, stop the consumer
			if ctx.Err() != nil {
				return
			}

			// Reset readiness so we wait for the next session
			consumer.ready = make(chan bool)
		}
	}()


	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			partition, offset := message.Partition, message.Offset
			log.Printf("Kafka websocket consumer: %s | Partition: %d | Offset: %d\n", message.Value, partition, offset)
			var websocketMessage models.WebsocketMessage
			err := json.Unmarshal(message.Value, &websocketMessage)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
			}
			websocketMessage.Origin = "server"
			hub.SendWebsocketMessage(websocketMessage)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}














// 	//////////////////
// 	group, err := sarama.NewConsumerGroup(brokers, "websocket-group", config)
// 	if err != nil {
// 		log.Fatalf("Failed to start Kafka consumer group: %v", err)
// 	}
// 	defer group.Close()

// 	fmt.Println("Kafka Consumer started...")

// 	for {
// 		select {
// 		case msg := <-partitionConsumer.Messages():
// 			var websocketMessage models.WebsocketMessage
// 			fmt.Printf("Kafka consumer viewed: %s\n", string(msg.Value))
// 			err := json.Unmarshal(msg.Value, &websocketMessage)
// 			if err != nil {
// 				fmt.Printf("Failed to unmarshal message: %v\n", err)
// 			}
// 			websocketMessage.Origin = "server"
// 			hub.SendWebsocketMessage(websocketMessage)
// 		case <-ctx.Done():
// 			fmt.Println("Consumer shutting down...")
// 			return
// 		}
// 	}
// }

// func CassandraConsumer(ctx context.Context, cassandraSession *gocql.Session) {

// }