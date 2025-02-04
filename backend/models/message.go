package models

type WebsocketMessage struct {
    Type string      `json:"type"` // message or circle
	Action string       `json:"action"` // create or delete
    Message *Message `json:"message,omitempty"` // optional
	Circle  *Circle  `json:"circle,omitempty"`
}

type Message struct {
	CircleID    string `json:"circle_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	AuthorID  string `json:"author_id"`
}
