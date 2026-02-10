package websocket

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	UserID         uuid.UUID
	conn           *websocket.Conn
	hub            *Hub
	send           chan *Message
	messageService services.MessageService
	roomService    services.RoomService
}

func NewClient(userID uuid.UUID, conn *websocket.Conn, hub *Hub, messageService services.MessageService, roomService services.RoomService) *Client {
	return &Client{
		UserID:         userID,
		conn:           conn,
		hub:            hub,
		send:           make(chan *Message, 256),
		messageService: messageService,
		roomService:    roomService,
	}
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		// Set sender ID from authenticated user
		message.SenderID = c.UserID

		// Save message to database
		ctx := context.Background()
		savedMsg, err := c.messageService.CreateMessage(ctx, services.CreateMessageRequest{
			RoomID:    message.RoomID,
			SenderID:  c.UserID,
			Content:   message.Content,
			Type:      message.Type,
			ReplyToID: message.ReplyToID,
		})

		if err != nil {
			log.Printf("error saving message: %v", err)
			continue
		}

		// Update message with saved data (ID, timestamps, sender info, etc.)
		message.ID = savedMsg.ID
		message.CreatedAt = savedMsg.CreatedAt
		message.Sender = &MessageSender{
			ID:       savedMsg.Sender.ID,
			Username: savedMsg.Sender.Username,
			FullName: savedMsg.Sender.FullName,
		}

		// Fetch room members to determine recipients
		room, err := c.roomService.GetRoomByID(ctx, message.RoomID)
		if err != nil {
			log.Printf("error fetching room members: %v", err)
			continue
		}

		recipients := make([]uuid.UUID, 0, len(room.Members))
		for _, member := range room.Members {
			recipients = append(recipients, member.UserID)
		}
		message.Recipients = recipients

		// Broadcast to all clients in the room
		c.hub.Broadcast <- &message
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling message: %v", err)
				continue
			}

			w.Write(messageBytes)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
