package routes

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/kevinsofyan/echoes-chat-api/docs"
	"github.com/kevinsofyan/echoes-chat-api/internal/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	AuthHandler      *handlers.AuthHandler
	UserHandler      *handlers.UserHandler
	WebSocketHandler *handlers.WebSocketHandler
	RoomHandler      *handlers.RoomHandler
}

func SetupRoutes(e *echo.Echo, h *Handlers) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	jwtSecret := os.Getenv("JWT_SECRET")
	log.Printf("JWT_SECRET loaded: %d characters", len(jwtSecret))

	jwtConfig := echojwt.Config{
		SigningKey: []byte(jwtSecret),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return jwt.MapClaims{}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			log.Printf("JWT Middleware Error: %v", err)
			log.Printf("Authorization Header: %s", c.Request().Header.Get("Authorization"))
			return echo.NewHTTPError(401, map[string]interface{}{
				"error":   "Invalid or expired token",
				"details": err.Error(),
			})
		},
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
		auth.POST("/register", h.AuthHandler.Register)
		auth.POST("/login", h.AuthHandler.Login)
		auth.POST("/logout", h.AuthHandler.Logout, echojwt.WithConfig(jwtConfig))
	}

	users := api.Group("/users")
	users.Use(echojwt.WithConfig(jwtConfig))
	{
		users.GET("/me", h.UserHandler.GetMe)
		users.GET("", h.UserHandler.GetAllUsers)
		users.GET("/:id", h.UserHandler.GetUserByID)
		users.PUT("/:id", h.UserHandler.UpdateUser)
		users.DELETE("/:id", h.UserHandler.DeleteUser)
	}

	// Room routes
	rooms := api.Group("/rooms")
	rooms.Use(echojwt.WithConfig(jwtConfig))
	{
		rooms.POST("", h.RoomHandler.CreateRoom)
		rooms.GET("/my", h.RoomHandler.GetMyRooms)
		rooms.GET("/:id", h.RoomHandler.GetRoomByID)
	}

	// WebSocket routes
	ws := api.Group("/ws")
	ws.Use(echojwt.WithConfig(jwtConfig))
	{
		ws.GET("/chat", h.WebSocketHandler.HandleWebSocket)
	}
}
