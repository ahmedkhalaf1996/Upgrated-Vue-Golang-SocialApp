package main

import (
	"Server/database"
	"log"

	_ "Server/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Fiber Golang Mongo Websocket Social App
// @version 1.0
// @description This is Swagger docs for rest api golang fiber with realtime chat
// @host localhost:5000
// @BasePath /
// @schemes http
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the token

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	database.Connect()

	// Create Fiber app
	app := fiber.New()

	// Setup CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Social App with Realtime Chat")
	})

	// Setup API routes
	SetupAPI(app)

	// Setup Realtime Chat
	SetupRealtimeChat(app)

	// Setup Realtime Notifications
	SetupRealtimeNotifications(app)

	// Serve swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Println("🚀 Server starting on port 5000 with API and Realtime Chat...")
	log.Fatal(app.Listen(":5000"))
}
