package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/websockets"
)

var brokers = []string{"localhost:9092"}
var topic = "chat"

var kafkaProducer sarama.SyncProducer

func WebsocketProducer(ctx context.Context, hub *websockets.Hub) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	kafkaProducer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka Websocket producer: %v", err)
	}
	defer kafkaProducer.Close()

	fmt.Println("Kafka Websocket Producer started...")

	for {
		select {
		case message := <-hub.Broadcast:
			var websocketMessage models.WebsocketMessage
			websocketMessage.Message = &models.Message{}
			websocketMessage.Circle = &models.Circle{}
			log.Printf("Client sent over websocket: %s\n", message)
			err := json.Unmarshal(message, &websocketMessage)
			if err != nil {
				log.Println("Error decoding JSON:", err)
				break
			}
			if websocketMessage.Origin == "server" {
				continue
			}
			switch websocketMessage.Type {
			case "message":
				if websocketMessage.Message != nil {
					log.Printf("New Message: %s from %s\n", websocketMessage.Message.Content, websocketMessage.Message.AuthorID)
					msg := &sarama.ProducerMessage{
						Topic: topic,
						Value: sarama.StringEncoder(message),
					}
					_, _, err = kafkaProducer.SendMessage(msg)
					if err != nil {
						fmt.Printf("Failed to send message: %v\n", err)
					} else {
						fmt.Println("Produced message successfully!")
					}
				}
			case "circle":
				if websocketMessage.Circle != nil {
					log.Printf("Circle %s %s\n", websocketMessage.Action, websocketMessage.Circle.Name)
				}
			default:
				log.Println("Unknown message type:", websocketMessage.Type)
			}
		}

	}
}

// TestProducer sends a message to Kafka every 5 seconds
func TestProducer(ctx context.Context) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	kafkaProducer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	fmt.Println("Kafka Producer started...")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Producer shutting down...")
			return
		default:
			message := models.Message{
				CircleID:  "1",
				Content:   "Hello from server",
				CreatedAt: time.Now().String(),
				AuthorID:  "1",
			}
			JsonMessage, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Failed to marshal message: %v\n", err)
			}

			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(JsonMessage),
			}

			_, _, err = kafkaProducer.SendMessage(msg)
			if err != nil {
				fmt.Printf("Failed to send message: %v\n", err)
			} else {
				fmt.Println("Produced message successfully!")
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func ChatProducer(websocketMessage models.WebsocketMessage) {

}
