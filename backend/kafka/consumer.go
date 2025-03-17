package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/IBM/sarama"
	"github.com/Leo7Deng/ChatApp/cassandra"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/Leo7Deng/ChatApp/websockets"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var hub *websockets.Hub
var cassandraSession *gocql.Session
var pool *pgxpool.Pool

func WebsocketConsumer(ctx context.Context, websocketHub *websockets.Hub) {
	hub = websocketHub
	groupID := "websocket-group"
	handler := &WebsocketConsumerHandler{hub: websocketHub}
	runConsumer(ctx , groupID, handler)
}

func CassandraConsumer(ctx context.Context, session *gocql.Session) {
	cassandraSession = session
	groupID := "cassandra-group"
	handler := &CassandraConsumerHandler{}
	runConsumer(ctx, groupID, handler)
}

func PostgresConsumer(ctx context.Context, postgresPool *pgxpool.Pool) {
	pool = postgresPool
	groupID := "postgres-group"
	handler := &PostgresConsumerHandler{}
	runConsumer(ctx, groupID, handler)
}


func runConsumer(ctx context.Context, groupID string, handler sarama.ConsumerGroupHandler) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Panicf("Error creating Kafka consumer group (%s): %v", groupID, err)
	}
	defer client.Close()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			err := client.Consume(ctx, []string{topic}, handler)
			if err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Printf("Error from Kafka consumer (%s): %v", groupID, err)
			}

			// If context is canceled, stop the consumer
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Printf("%s consumer started...", groupID)
	<-ctx.Done()
	log.Printf("%s consumer shutting down...", groupID)
	wg.Wait()
}


type WebsocketConsumerHandler struct {hub *websockets.Hub}
type CassandraConsumerHandler struct {}
type PostgresConsumerHandler struct {}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *WebsocketConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumer *CassandraConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumer *PostgresConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *WebsocketConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumer *CassandraConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumer *PostgresConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *WebsocketConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
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
			log.Printf("Kafka consumer: %s | Partition: %d | Offset: %d\n", message.Value, partition, offset)
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

func (consumer *CassandraConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			partition, offset := message.Partition, message.Offset
			var websocketMessage models.WebsocketMessage
			err := json.Unmarshal(message.Value, &websocketMessage)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
			}
			log.Printf("Cassandra consumer: %s | Partition: %d | Offset: %d\n", message.Value, partition, offset)
			insertMessage := models.Message{
				CircleID: websocketMessage.Message.CircleID,
				AuthorID: websocketMessage.Message.AuthorID,
				Content:  websocketMessage.Message.Content,
				CreatedAt: websocketMessage.Message.CreatedAt,
			}
			err = cassandra.InsertMessage(cassandraSession, insertMessage)

			// Handle error on unprocessed message insert into Cassandra
			if err != nil {
				fmt.Printf("Failed to insert message: %v\n", err)
			} else {
				log.Printf("Message inserted into Cassandra: %v\n", insertMessage)
				session.MarkMessage(message, "")
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func (consumer *PostgresConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			partition, offset := message.Partition, message.Offset
			var websocketMessage models.WebsocketMessage
			err := json.Unmarshal(message.Value, &websocketMessage)
			if err != nil {
				fmt.Printf("Failed to unmarshal message: %v\n", err)
			}
			log.Printf("Postgres consumer: %s | Partition: %d | Offset: %d\n", message.Value, partition, offset)
			type InsertMessage struct {
				MessageID string `json:"message_id"`
				Message models.Message `json:"message"`
			}
			insertMessage := InsertMessage{
				MessageID: uuid.New().String(),
				Message: models.Message{
					CircleID: websocketMessage.Message.CircleID,
					AuthorID: websocketMessage.Message.AuthorID,
					Content:  websocketMessage.Message.Content,
					CreatedAt: websocketMessage.Message.CreatedAt,
				},
			}
			conn, err := pool.Acquire(session.Context())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to acquire a connection from the pool: %v\n", err)
				return err
			}
			defer conn.Release()

			_, err = conn.Exec(
				session.Context(),
				"INSERT INTO messages (message_id, circle_id, author_id, content) VALUES ($1, $2, $3, $4)",
				insertMessage.MessageID,
				insertMessage.Message.CircleID,
				insertMessage.Message.AuthorID,
				insertMessage.Message.Content,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error inserting into messages: %v\n", err)
				return err
			}
		case <-session.Context().Done():
			return nil
		}
	}
}