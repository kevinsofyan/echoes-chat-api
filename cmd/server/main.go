package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kevinsofyan/echoes-chat-api/internal/container"
	"github.com/kevinsofyan/echoes-chat-api/internal/database"
	"github.com/kevinsofyan/echoes-chat-api/internal/routes"
	"github.com/labstack/echo/v4"
)

// @title Echoes Chat API
// @version 1.0
// @description This is a chat application API with WebSocket support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@echoes-chat.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	c := container.NewContainer(db)
	e := echo.New()
	go c.Hub.Run()
	routes.SetupRoutes(e, c.Handlers)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
