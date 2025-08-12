package main

import (
	"Server/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

// SetupAPI configures all REST API routes
func SetupAPI(app *fiber.App) {
	log.Println("Setting up API routes...")

	// Setup all existing routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupPostRoutes(app)
	routes.SetupChatRoutes(app)
	routes.SetupNotificationRoutes(app)

	// API health check
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"service":  "api",
			"database": "connected",
		})
	})

	log.Println("✅ API routes configured successfully")
}
