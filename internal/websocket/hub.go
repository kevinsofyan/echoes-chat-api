package websocket

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	clients map[uuid.UUID]*Client

	Broadcast chan *Message

	Register chan *Client

	Unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	ID        uuid.UUID  `json:"id,omitempty"`
	RoomID    uuid.UUID  `json:"room_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Content   string     `json:"content"`
	Type      string     `json:"type"` // "text", "image", "file", "video", "audio"
	FileURL   string     `json:"file_url,omitempty"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
	CreatedAt time.Time  `json:"created_at,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.UserID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID uuid.UUID, message *Message) {
	h.Broadcast <- message
}
