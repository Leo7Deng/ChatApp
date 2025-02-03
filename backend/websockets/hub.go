package websockets

import (
	"encoding/json"
	"log"

	"github.com/Leo7Deng/ChatApp/models"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// map of userID to clients (could have multiple clients per user)
	userClients map[string]map[*Client]bool

	// map of circleID to userIDs (could have multiple users per circle)
	circleUsers map[string]map[string]bool
}

func NewHub() *Hub {
	return &Hub{
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		userClients: make(map[string]map[*Client]bool),
		circleUsers: make(map[string]map[string]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			if _, ok := h.userClients[client.userID]; !ok {
				h.userClients[client.userID] = make(map[*Client]bool)
			}
			h.userClients[client.userID][client] = true
			for _, circle := range client.circles {
				if _, ok := h.circleUsers[circle]; !ok {
					h.circleUsers[circle] = make(map[string]bool)
				}
				h.circleUsers[circle][client.userID] = true
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				if connections, ok := h.userClients[client.userID]; ok {
					delete(connections, client)
					if len(connections) == 0 {
						delete(h.userClients, client.userID)
					}
				}
				for _, circle := range client.circles {
					if users, ok := h.circleUsers[circle]; ok {
						delete(users, client.userID)
						if len(users) == 0 {
							delete(h.circleUsers, circle)
						}
					}
				}
			}
		case message := <-h.broadcast:
			log.Printf("Websocket message: %s\n", message)
			var msg models.Message
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v\n", err)
			}
			for userID := range h.circleUsers[msg.CircleID] {
				for client := range h.userClients[userID] {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

func (h *Hub) SendToUser(userID string, message []byte) {
	for client := range h.clients {
		client.hub.broadcast <- message
	}
}

func (h *Hub) AddUsersToCircle(circleID string, userID []string) {
	if _, ok := h.circleUsers[circleID]; !ok {
		h.circleUsers[circleID] = make(map[string]bool)
	}
	for _, id := range userID {
		h.circleUsers[circleID][id] = true
	}
	log.Printf("Users added to circle %s: %v\n", circleID, userID)
	log.Printf("Circle users: %v\n", h.circleUsers)
}