package models

import "time"

type Circle struct {
	ID        string       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCircleData struct {
	Name string `json:"name"`
}