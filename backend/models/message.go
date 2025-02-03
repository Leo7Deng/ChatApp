package models

type Message struct {
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
	CircleID    string `json:"chat_id"`
}
