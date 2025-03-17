package models

import "time"

type WebsocketMessage struct {
	Origin  string      `json:"origin"` // user or server
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

type SearchMessage struct {
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	AuthorUsername string `json:"author_username"`
}