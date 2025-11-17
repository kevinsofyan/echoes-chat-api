package container

import (
	"github.com/kevinsofyan/echoes-chat-api/internal/handlers"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
	"github.com/kevinsofyan/echoes-chat-api/internal/routes"
	"github.com/kevinsofyan/echoes-chat-api/internal/services"
	"gorm.io/gorm"
)

type Container struct {
	Handlers *routes.Handlers
}

func NewContainer(db *gorm.DB) *Container {
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	allHandlers := &routes.Handlers{
		UserHandler: userHandler,
	}

	return &Container{
		Handlers: allHandlers,
	}
}
