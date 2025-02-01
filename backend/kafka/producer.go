package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

var brokers = []string{"localhost:9092"}
var topic = "chat"

// StartProducer sends a message to Kafka every 5 seconds
func StartProducer(ctx context.Context) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer producer.Close()

	fmt.Println("Kafka Producer started...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Producer shutting down...")
			return
		default:
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(fmt.Sprintf("Hello from server at %s", time.Now().Format(time.RFC3339))),
			}

			_, _, err := producer.SendMessage(msg)
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			} else {
				fmt.Println("Produced message successfully!")
			}

			time.Sleep(5 * time.Second)
		}
	}
}