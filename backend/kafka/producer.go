package kafka

import (
	"context"
	"fmt"
	"log"
	"time"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/Leo7Deng/ChatApp/models"
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
			message := models.Message{
				Content: "Hello from server",
				CreatedAt: time.Now().String(),
				AuthorID: "1",
				CircleID: "1",
			}
			JsonMessage, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Failed to marshal message: %v\n", err)
			}

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(JsonMessage),
			}

			_, _, err = producer.SendMessage(msg)
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			} else {
				fmt.Println("Produced message successfully!")
			}

			time.Sleep(5 * time.Second)
		}
	}
}