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

type MessageSender struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
}

type Message struct {
	ID         uuid.UUID      `json:"id,omitempty"`
	RoomID     uuid.UUID      `json:"room_id"`
	SenderID   uuid.UUID      `json:"sender_id,omitempty"`
	Sender     *MessageSender `json:"sender,omitempty"`
	Content    string         `json:"content"`
	Type       string         `json:"type"` // "text", "image", "file", "video", "audio"
	FileURL    string         `json:"file_url,omitempty"`
	ReplyToID  *uuid.UUID     `json:"reply_to_id,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
	Recipients []uuid.UUID    `json:"-"` // List of user IDs to receive the message
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
			// If recipients are specified, only send to them
			if len(message.Recipients) > 0 {
				for _, recipientID := range message.Recipients {
					if client, ok := h.clients[recipientID]; ok {
						select {
						case client.send <- message:
						default:
							close(client.send)
							h.mu.RUnlock() // Unlock before modifying map
							h.mu.Lock()
							delete(h.clients, client.UserID)
							h.mu.Unlock()
							h.mu.RLock() // Re-lock for reading
						}
					}
				}
			} else {
				// Fallback to broadcast all (or maybe we should disable this?)
				// For now, let's keep it but it shouldn't be used for chat messages
				for _, client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						h.mu.RUnlock()
						h.mu.Lock()
						delete(h.clients, client.UserID)
						h.mu.Unlock()
						h.mu.RLock()
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(roomID uuid.UUID, message *Message) {
	h.Broadcast <- message
}
