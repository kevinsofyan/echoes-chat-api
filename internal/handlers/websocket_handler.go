package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
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
	roomService    services.RoomService
}

func NewWebSocketHandler(hub *ws.Hub, messageService services.MessageService, roomService services.RoomService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:            hub,
		messageService: messageService,
		roomService:    roomService,
	}
}

// HandleWebSocket godoc
// @Summary WebSocket endpoint for real-time chat
// @Tags websocket
// @Param token query string true "JWT Token"
// @Router /api/v1/ws/chat [get]
func (h *WebSocketHandler) HandleWebSocket(c echo.Context) error {
	// Get token from query parameter (for browser WebSocket compatibility)
	tokenString := c.QueryParam("token")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "token required in query parameter",
		})
	}

	// Validate JWT token
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid or expired token",
		})
	}

	// Extract user ID from token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid token claims",
		})
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "user_id not found in token",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid user_id format",
		})
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return err
	}

	client := ws.NewClient(userID, conn, h.hub, h.messageService, h.roomService)

	h.hub.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()

	return nil
}
