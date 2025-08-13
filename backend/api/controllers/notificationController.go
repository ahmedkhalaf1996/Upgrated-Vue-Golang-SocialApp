package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NotificationResponse represents the response structure with populated user data
type NotificationResponse struct {
	ID        primitive.ObjectID `json:"_id,omitempty"`
	Deatils   string             `json:"deatils"`
	MainUID   string             `json:"mainuid"`
	TargetID  string             `json:"targetid"`
	IsReaded  bool               `json:"isreded"`
	CreatedAt time.Time          `json:"createdAt"`
	User      UserData           `json:"user"`
}

type UserData struct {
	Name     string `json:"name"`
	ImageUrl string `json:"imageUrl"`
}

func populateNotificationsWithUserData(ctx context.Context, notifications []models.Notification) ([]NotificationResponse, error) {
	var response []NotificationResponse
	userCol := database.DB.Collection("users")

	for _, notification := range notifications {
		// Convert UserID string to ObjectID
		userObjID, err := primitive.ObjectIDFromHex(notification.UserID)
		if err != nil {
			continue // Skip invalid user IDs
		}

		// Fetch fresh user data
		var user models.UserModel
		err = userCol.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
		if err != nil {
			// If user not found, use default values
			user.Name = "Unknown User"
			user.ImageUrl = ""
		}

		// Create response with fresh user data
		notificationResp := NotificationResponse{
			ID:        notification.ID,
			Deatils:   notification.Deatils,
			MainUID:   notification.MainUID,
			TargetID:  notification.TargetID,
			IsReaded:  notification.IsReaded,
			CreatedAt: notification.CreatedAt,
			User: UserData{
				Name:     user.Name,
				ImageUrl: user.ImageUrl,
			},
		}
		response = append(response, notificationResp)
	}

	return response, nil
}

func MarknotAsReaded(c *fiber.Ctx) error {
	// parse query parameter
	id := c.Query("id")
	if id == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "id in query is Required",
		})
	}

	// construct the filter and update
	filter := bson.M{"mainuid": bson.M{"$regex": id, "$options": "i"}}
	update := bson.M{"$set": bson.M{"isreded": true}}

	// update
	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := NotificationSchema.UpdateMany(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to mark notifications as read",
			"error":   err.Error(),
		})
	}

	// retrieve the updated notifications
	options := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := NotificationSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve the updated notifications",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to decode notifications",
			"error":   err.Error(),
		})
	}

	// Populate with fresh user data
	response, err := populateNotificationsWithUserData(ctx, notifications)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to populate user data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": response,
	})
}

func GetUserNotification(c *fiber.Ctx) error {
	// parse query parameter
	id := c.Params("userid")
	if id == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "id in Params is Required",
		})
	}

	var NotificationSchema = database.DB.Collection("notifications")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// construct the filter
	filter := bson.M{"mainuid": bson.M{"$regex": id, "$options": "i"}}
	options := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	// retrieve notifications
	cursor, err := NotificationSchema.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve notifications",
			"error":   err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to decode notifications",
			"error":   err.Error(),
		})
	}

	if len(notifications) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"notifications": []NotificationResponse{},
		})
	}

	// Populate with fresh user data
	response, err := populateNotificationsWithUserData(ctx, notifications)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to populate user data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": response,
	})
}
