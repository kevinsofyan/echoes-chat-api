package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"github.com/kevinsofyan/echoes-chat-api/internal/utils"
	ws "github.com/kevinsofyan/echoes-chat-api/internal/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now (change in production)
		return true
	},
}

type WebSocketHandler struct {
	hub            *ws.Hub
	messageService services.MessageService
}

func NewWebSocketHandler(hub *ws.Hub, messageService services.MessageService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:            hub,
		messageService: messageService,
	}
}

// HandleWebSocket godoc
// @Summary WebSocket endpoint for real-time chat
// @Tags websocket
// @Security BearerAuth
// @Router /api/v1/ws/chat [get]
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// Get user ID from JWT token
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized",
		})
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return err
	}

	client := ws.NewClient(userID, conn, h.hub, h.messageService)

	h.hub.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	return nil
}
