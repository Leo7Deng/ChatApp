package kafka

import (
	"context"
	"sync"
	"github.com/Leo7Deng/ChatApp/websockets"
	"github.com/gocql/gocql"
)

// InitKafka starts the Kafka producer and consumer as goroutines
func InitKafka(ctx context.Context, wg *sync.WaitGroup, hub *websockets.Hub, cassandraSession *gocql.Session) {
	wg.Add(2)

	go func() {
		defer wg.Done()
		WebsocketConsumer(ctx, hub) 
	}()

	go func() {
		defer wg.Done()
		CassandraConsumer(ctx, cassandraSession)
	}()

	go func() {
		defer wg.Done()
		WebsocketProducer(ctx, hub)
	}()
}