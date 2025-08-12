package realtime

import (
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

type Notification struct {
	ID        string    `json:"_id"`
	Details   string    `json:"details"`
	MainUID   string    `json:"mainuid"`
	TargetID  string    `json:"targetid"`
	IsReaded  bool      `json:"isreded"`
	CreatedAt time.Time `json:"createdAt"`
	User      User      `json:"user"`
}

type User struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type NotificationManager struct {
	connections map[string]*websocket.Conn
	lock        sync.RWMutex
}

// Global notification manager instance
var notificationManager *NotificationManager

func init() {
	notificationManager = &NotificationManager{
		connections: make(map[string]*websocket.Conn),
	}
}

// GetNotificationManager returns the global notification manager
func GetNotificationManager() *NotificationManager {
	return notificationManager
}

// AddNotificationConnection adds a WebSocket connection for notifications
func (nm *NotificationManager) AddNotificationConnection(userID string, conn *websocket.Conn) {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	nm.connections[userID] = conn
	log.Printf("📢 User %s connected to notifications", userID)
}

// RemoveNotificationConnection removes a WebSocket connection
func (nm *NotificationManager) RemoveNotificationConnection(userID string) {
	nm.lock.Lock()
	defer nm.lock.Unlock()

	delete(nm.connections, userID)
	log.Printf("📢 User %s disconnected from notifications", userID)
}

// SendNotificationToUser sends a notification to a specific user
func (nm *NotificationManager) SendNotificationToUser(userID string, notification Notification) error {
	nm.lock.RLock()
	conn, exists := nm.connections[userID]
	nm.lock.RUnlock()

	if !exists {
		log.Printf("📢 User %s not connected for notifications", userID)
		return nil // Not an error, user just isn't online
	}

	err := conn.WriteJSON(notification)
	if err != nil {
		log.Printf("📢 Error sending notification to user %s: %v", userID, err)
		// Remove the connection if it's broken
		nm.RemoveNotificationConnection(userID)
		return err
	}

	log.Printf("📢 Notification sent to user %s: %s", userID, notification.Details)
	return nil
}

// GetConnectedUsers returns list of connected user IDs
func (nm *NotificationManager) GetConnectedUsers() []string {
	nm.lock.RLock()
	defer nm.lock.RUnlock()

	users := make([]string, 0, len(nm.connections))
	for userID := range nm.connections {
		users = append(users, userID)
	}
	return users
}
