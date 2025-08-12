package main

import (
	"Server/realtime"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRealtimeNotifications configures WebSocket routes for notifications
func SetupRealtimeNotifications(app *fiber.App) {
	log.Println("Setting up Realtime Notifications...")

	// Get the notification manager
	notificationManager := realtime.GetNotificationManager()

	// Notification status endpoint
	app.Get("/notifications/status", func(c *fiber.Ctx) error {
		connectedUsers := notificationManager.GetConnectedUsers()
		return c.JSON(fiber.Map{
			"status":          "active",
			"service":         "realtime-notifications",
			"websocket":       "/notifications/ws/{userId}",
			"connected_users": len(connectedUsers),
			"users":           connectedUsers,
		})
	})

	// WebSocket route for notifications
	app.Get("/notifications/ws/:userId", websocket.New(func(c *websocket.Conn) {
		userID := c.Params("userId")

		log.Printf("📢 New notification WebSocket connection for user: %s", userID)

		// Add connection to notification manager
		notificationManager.AddNotificationConnection(userID, c)

		defer func() {
			log.Printf("📢 Closing notification WebSocket connection for user: %s", userID)
			notificationManager.RemoveNotificationConnection(userID)
			c.Close()
		}()

		// Keep connection alive and handle incoming messages (if any)
		for {
			// For notifications, we mainly send data to client, but we can still listen
			// for ping/pong or any client messages to keep connection alive
			_, _, err := c.ReadMessage()
			if err != nil {
				log.Printf("📢 Notification WebSocket error for user %s: %v", userID, err)
				break
			}
		}
	}))

	log.Println("✅ Realtime Notifications configured successfully")
	log.Println("📢 Notification WebSocket endpoint available at: ws://localhost:5000/notifications/ws/{userId}")
}
