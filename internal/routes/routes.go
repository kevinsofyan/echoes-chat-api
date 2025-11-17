package routes

import (
	"os"

	_ "github.com/kevinsofyan/echoes-chat-api/docs"
	"github.com/kevinsofyan/echoes-chat-api/internal/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	UserHandler *handlers.UserHandler
}

func SetupRoutes(e *echo.Echo, h *Handlers) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	jwtConfig := echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")

	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.UserHandler.Register)
		auth.POST("/login", h.UserHandler.Login)
	}

	users := api.Group("/users")
	users.Use(echojwt.WithConfig(jwtConfig))
	{
		users.GET("/me", h.UserHandler.GetMe) // Get current user
		users.GET("", h.UserHandler.GetAllUsers)
		users.GET("/:id", h.UserHandler.GetUserByID)
		users.PUT("/:id", h.UserHandler.UpdateUser)    // Only allow updating own profile
		users.DELETE("/:id", h.UserHandler.DeleteUser) // Only allow deleting own account
	}
}
