package cassandra

import (
	"fmt"
	"time"

	"github.com/Leo7Deng/ChatApp/models"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func InsertMessage(session *gocql.Session, message models.Message) error {
	messageID := uuid.New().String()
	// convert time string to cassandra timestamp
	parsedTime, err := time.Parse(time.RFC3339, message.CreatedAt)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return err
	}
	query := session.Query(`
		INSERT INTO chat_app.messages (circle_id, created_at, message_id, author_id, content)
		VALUES (?, ?, ?, ?, ?)`,
		message.CircleID, parsedTime, messageID, message.AuthorID, message.Content)
	err = query.Exec()
	if err != nil {
		fmt.Println("Error inserting message into Cassandra:", err)
		return err
	}
	return nil
}