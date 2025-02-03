package kafka

import (
	"context"
	"fmt"
	"log"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/Leo7Deng/ChatApp/models"
)

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