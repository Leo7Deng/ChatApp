package cassandra

import (
	"fmt"
	"github.com/Leo7Deng/ChatApp/models"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
)

func InsertMessage(session *gocql.Session, message models.Message) error {
	messageID := uuid.New().String()
	query := session.Query(`
		INSERT INTO chat_app.messages (circle_id, created_at, message_id, author_id, content)
		VALUES (?, ?, ?, ?, ?)`,
		message.CircleID, message.CreatedAt, messageID, message.AuthorID, message.Content)
	err := query.Exec()
	if err != nil {
		fmt.Println("Error inserting message into Cassandra:", err)
		return err
	}
	return nil
}