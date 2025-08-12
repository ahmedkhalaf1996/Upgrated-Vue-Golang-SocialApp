package services

import (
	"Server/models"
	"Server/realtime"
	"log"
)

// SendNotification sends a notification via WebSocket (replaces gRPC call)
func SendNotification(notification models.Notification) error {
	// Convert models.Notification to realtime.Notification
	realtimeNotification := realtime.Notification{
		ID:        notification.ID.Hex(),
		Details:   notification.Deatils, // Note: your model has a typo "Deatils"
		MainUID:   notification.MainUID,
		TargetID:  notification.TargetID,
		IsReaded:  notification.IsReaded,
		CreatedAt: notification.CreatedAt,
		User: realtime.User{
			Name:   notification.User.Name,
			Avatar: notification.User.Avatart, // Note: your model has a typo "Avatart"
		},
	}

	// Get the notification manager and send to the target user
	notificationManager := realtime.GetNotificationManager()
	err := notificationManager.SendNotificationToUser(notification.MainUID, realtimeNotification)

	if err != nil {
		log.Printf("Failed to send realtime notification: %v", err)
		return err
	}

	log.Printf("✅ Notification sent successfully to user %s", notification.MainUID)
	return nil
}
