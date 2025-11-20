package container

import (
	"github.com/kevinsofyan/echoes-chat-api/internal/handlers"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
	"github.com/kevinsofyan/echoes-chat-api/internal/routes"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"github.com/kevinsofyan/echoes-chat-api/internal/websocket"
	"gorm.io/gorm"
)

type Container struct {
	Handlers *routes.Handlers
	Hub      *websocket.Hub
}

func NewContainer(db *gorm.DB) *Container {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	roomRepo := repositories.NewRoomRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenRepo)
	userService := services.NewUserService(userRepo)
	messageService := services.NewMessageService(messageRepo)
	roomService := services.NewRoomService(roomRepo)

	// Initialize WebSocket hub
	hub := websocket.NewHub()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	wsHandler := handlers.NewWebSocketHandler(hub, messageService)
	roomHandler := handlers.NewRoomHandler(roomService)

	// Group handlers
	allHandlers := &routes.Handlers{
		AuthHandler:      authHandler,
		UserHandler:      userHandler,
		WebSocketHandler: wsHandler,
		RoomHandler:      roomHandler,
	}

	return &Container{
		Handlers: allHandlers,
		Hub:      hub,
	}
}
