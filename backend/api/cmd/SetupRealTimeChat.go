package main

import (
	"Server/realtime"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRealtimeChat configures all realtime chat WebSocket routes and functionality
func SetupRealtimeChat(app *fiber.App) {
	log.Println("Setting up Realtime Chat...")

	// Initialize connection manager with user friends function
	manager := realtime.NewConnectionManager(realtime.GetUserFriends)

	// Realtime chat status endpoint
	app.Get("/realtime/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":      "active",
			"service":     "realtime-chat",
			"websocket":   "/ws/{userId}",
			"connections": "active", // You could expose actual count if needed
		})
	})

	// WebSocket route for realtime chat
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		if manager == nil {
			log.Printf("Error: Connection manager is nil for user %s", id)
			return
		}

		log.Printf("New WebSocket connection for user: %s", id)

		// Add connection to manager
		manager.AddConnection(id, c)
		defer func() {
			log.Printf("Closing WebSocket connection for user: %s", id)
			manager.RemoveConnection(id)
			c.Close()
		}()

		// Handle incoming messages
		var msg realtime.Message
		for {
			err := c.ReadJSON(&msg)
			if err != nil {
				handleWebSocketError(err, id)
				manager.RemoveConnection(id)
				c.Close()
				break
			}

			log.Printf("📨 Received message from %s to %s: %s", msg.Sender, msg.Recever, msg.Content)

			// Send message through manager
			manager.SendToReceiver(msg)
		}
	}))

	log.Println("✅ Realtime Chat configured successfully")
	log.Println("💬 WebSocket endpoint available at: ws://localhost:5000/ws/{userId}")
}

// handleWebSocketError handles WebSocket connection errors
func handleWebSocketError(err error, userID string) {
	log.Printf("❌ WebSocket error for user %s: %v", userID, err)
}
