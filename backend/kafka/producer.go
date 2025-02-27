package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/websockets"
)

var topic = "chat"
var partitions = int32(3)
var brokers = []string{"kafka:9093"}
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

	for message := range hub.Broadcast {
		var websocketMessage models.WebsocketMessage
		websocketMessage.Message = &models.Message{}
		websocketMessage.Circle = &models.Circle{}
		log.Printf("Client sent over websocket: %s\n", message)
		err := json.Unmarshal(message, &websocketMessage)
		if (err != nil) {
			log.Println("Error decoding JSON:", err)
			continue
		}
		if websocketMessage.Origin == "server" {
			continue
		}
		switch websocketMessage.Type {
		case "message":
			if websocketMessage.Message != nil {
				// log.Printf("New Message: %s from %s\n", websocketMessage.Message.Content, websocketMessage.Message.AuthorID)
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(message),
				}
				partition, offset, err := kafkaProducer.SendMessage(msg)
				if err != nil {
					fmt.Printf("Failed to send message: %v\n", err)
				} else {
					fmt.Printf("Websocket producer sent message to partition %d at offset %d\n", partition, offset)
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
