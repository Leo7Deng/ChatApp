package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
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
			fmt.Printf("Received message: %s\n", string(msg.Value))
		case <-ctx.Done():
			fmt.Println("Consumer shutting down...")
			return
		}
	}
}