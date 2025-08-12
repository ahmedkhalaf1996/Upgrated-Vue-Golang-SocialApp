package realtime

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
)

// StartRealtimeServer starts the realtime chat server
func StartRealtimeServer(port string) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	manager := NewConnectionManager(GetUserFriends)

	// Register WebSocket route
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		if manager == nil {
			return
		}
		manager.AddConnection(id, c)
		defer func() {
			manager.RemoveConnection(id)
			c.Close()
		}()

		var msg Message
		for {
			err := c.ReadJSON(&msg)
			if err != nil {
				handleWebSocketError(err, id)
				manager.RemoveConnection(id)
				c.Close()
				break
			}

			log.Printf("Received message from %s to %s: %s", msg.Sender, msg.Recever, msg.Content)
			manager.SendToReceiver(msg)
		}
	}))

	log.Printf("Realtime chat server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func handleWebSocketError(err error, userID string) {
	log.Printf("WebSocket error for user %s: %v", userID, err)
}
