package services

import (
	"Server/database"
	"Server/models"
	"Server/realtime"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendNotification sends a notification via WebSocket with fresh user data
func SendNotification(notification models.Notification) error {
	// Fetch fresh user data for real-time notification
	userCol := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(notification.UserID)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return err
	}

	var user models.UserModel
	err = userCol.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		log.Printf("Failed to fetch user data: %v", err)
		// Use default values if user not found
		user.Name = "Unknown User"
		user.ImageUrl = ""
	}

	// Convert to realtime notification with fresh user data
	realtimeNotification := realtime.Notification{
		ID:        notification.ID.Hex(),
		Details:   notification.Deatils,
		MainUID:   notification.MainUID,
		TargetID:  notification.TargetID,
		IsReaded:  notification.IsReaded,
		CreatedAt: notification.CreatedAt,
		User: realtime.User{
			Name:     user.Name,
			ImageUrl: user.ImageUrl,
		},
	}

	// Get the notification manager and send to the target user
	notificationManager := realtime.GetNotificationManager()
	err = notificationManager.SendNotificationToUser(notification.MainUID, realtimeNotification)

	if err != nil {
		log.Printf("Failed to send realtime notification: %v", err)
		return err
	}

	log.Printf("✅ Notification sent successfully to user %s", notification.MainUID)
	return nil
}
