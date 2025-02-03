package models

type Message struct {
	CircleID    string `json:"chat_id"`
	Type        string `json:"type"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
}
