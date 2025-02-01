package kafka

import (
	"context"
	"sync"
)

// InitKafka starts the Kafka producer and consumer as goroutines
func InitKafka(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(2)

	go func() {
		defer wg.Done()
		StartConsumer(ctx)
	}()

	go func() {
		defer wg.Done()
		StartProducer(ctx)
	}()
}