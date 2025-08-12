package realtime

import (
	"Server/database"
	"Server/models"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserFriends(userID string) <-chan []string {
	ch := make(chan []string)
	go func() {
		defer close(ch)

		// Create a context with timeout for database operations
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userFriends, err := getUserFollowingFollowersFromDB(ctx, userID)
		if err != nil {
			log.Printf("Error getting friends for user %s: %v", userID, err)
			ch <- []string{}
			return
		}

		ch <- userFriends
	}()
	return ch
}

func getUserFollowingFollowersFromDB(ctx context.Context, userID string) ([]string, error) {
	UserSchema := database.DB.Collection("users")

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	var user models.UserModel

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %v", err)
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": uid}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Create a map to avoid duplicates
	friendsMap := make(map[string]bool)

	// Add following users
	for _, id := range user.Following {
		friendsMap[id] = true
	}

	// Add followers
	for _, id := range user.Followers {
		friendsMap[id] = true
	}

	// Convert map to slice
	var friends []string
	for id := range friendsMap {
		friends = append(friends, id)
	}

	log.Printf("Found %d friends for user %s", len(friends), userID)
	return friends, nil
}
