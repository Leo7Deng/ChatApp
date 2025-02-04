package kafka

import (
	"context"
	"sync"
	"github.com/Leo7Deng/ChatApp/websockets"
)

// InitKafka starts the Kafka producer and consumer as goroutines
func InitKafka(ctx context.Context, wg *sync.WaitGroup, hub *websockets.Hub) {
	wg.Add(3)

	go func() {
		defer wg.Done()
		StartConsumer(ctx) // Test consumer
	}()

	go func() {
		defer wg.Done()
		WebsocketConsumer(ctx, hub) // Websocket consumer
	}()

	go func() {
		defer wg.Done()
		WebsocketProducer(ctx, hub)
	}()
}