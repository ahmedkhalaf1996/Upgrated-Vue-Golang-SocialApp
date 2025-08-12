package realtime

import (
	"Server/database" // Update this import path to match your project structure
	"Server/models"   // Update this import path to match your project structure
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Sender  string `json:"sender"`
	Recever string `json:"recever"`
	Content string `json:"content"`
}

type ConnectionManager struct {
	connections    map[string]*websocket.Conn
	onlineFriends  map[string][]string
	getUserFriends func(string) <-chan []string
	lock           sync.Mutex
}

func NewConnectionManager(getUserFriends func(string) <-chan []string) *ConnectionManager {
	if getUserFriends != nil {
		return &ConnectionManager{
			connections:    make(map[string]*websocket.Conn),
			onlineFriends:  make(map[string][]string),
			getUserFriends: getUserFriends,
		}
	}
	return nil
}

func (cm *ConnectionManager) AddConnection(userID string, conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.connections[userID] = conn
	cm.onlineFriends[userID] = []string{}

	// Notify existing users about the new online friend
	for friendID := range cm.onlineFriends {
		if friendID != userID && cm.isFriend(userID, friendID) {
			cm.onlineFriends[friendID] = append(cm.onlineFriends[friendID], userID)
			err := cm.connections[friendID].WriteJSON(map[string]interface{}{
				"onlineFriends": cm.onlineFriends[friendID],
			})
			if err != nil {
				log.Printf("Error notifying %s about %s: %v", friendID, userID, err)
				return
			}
		}
	}

	// Update the online friends list for the new user
	go func() {
		for friends := range cm.getUserFriends(userID) {
			if friends == nil {
				continue
			}

			for _, friendID := range friends {
				if cm.connections[friendID] != nil {
					cm.onlineFriends[userID] = append(cm.onlineFriends[userID], friendID)
					err := cm.connections[userID].WriteJSON(map[string]interface{}{
						"onlineFriends": cm.onlineFriends[userID],
					})
					if err != nil {
						log.Printf("Error notifying %s about %s: %v", userID, friendID, err)
						return
					}
				}
			}
		}
	}()
}

func (cm *ConnectionManager) RemoveConnection(userID string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	delete(cm.connections, userID)
	delete(cm.onlineFriends, userID)
	for friendID := range cm.onlineFriends {
		for i, id := range cm.onlineFriends[friendID] {
			if id == userID {
				cm.onlineFriends[friendID] = append(cm.onlineFriends[friendID][:i], cm.onlineFriends[friendID][i+1:]...)
				err := cm.connections[friendID].WriteJSON(map[string]interface{}{
					"onlineFriends": cm.onlineFriends[friendID],
				})
				if err != nil {
					log.Printf("Error notifying %s about %s: %v", friendID, userID, err)
				}
				break
			}
		}
	}
}

func (cm *ConnectionManager) SendToReceiver(msg Message) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if conn, ok := cm.connections[msg.Recever]; ok {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to %s: %v", msg.Recever, err)
		}

		// Save Message To DB directly (no gRPC needed)
		err = cm.saveMessageToDB(msg)
		if err != nil {
			log.Printf("Error saving message to database: %v", err)
		}
	} else {
		log.Printf("Receiver %s not found", msg.Recever)
		// Still save the message even if receiver is offline
		err := cm.saveMessageToDB(msg)
		if err != nil {
			log.Printf("Error saving offline message to database: %v", err)
		}
	}
}

// saveMessageToDB directly saves the message to database
func (cm *ConnectionManager) saveMessageToDB(msg Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Validate sender and receiver exist
	if msg.Sender == "" || msg.Recever == "" {
		return fmt.Errorf("sender and receiver are required")
	}

	UserSchema := database.DB.Collection("users")

	// Check if sender and receiver exist
	senderID, err := primitive.ObjectIDFromHex(msg.Sender)
	if err != nil {
		return fmt.Errorf("invalid sender ID: %v", err)
	}

	receiverID, err := primitive.ObjectIDFromHex(msg.Recever)
	if err != nil {
		return fmt.Errorf("invalid receiver ID: %v", err)
	}

	var sender, receiver models.UserModel
	err = UserSchema.FindOne(ctx, bson.M{"_id": senderID}).Decode(&sender)
	if err != nil {
		return fmt.Errorf("sender not found: %v", err)
	}

	err = UserSchema.FindOne(ctx, bson.M{"_id": receiverID}).Decode(&receiver)
	if err != nil {
		return fmt.Errorf("receiver not found: %v", err)
	}

	// Save message
	message := models.Message{
		Content: msg.Content,
		Sender:  msg.Sender,
		Recever: msg.Recever,
	}

	_, err = database.DB.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to save message to db: %v", err)
	}

	// Update unread messages
	unreadedmessagesSchema := database.DB.Collection("unReadedmessages")

	existingRecord := bson.M{}
	err = unreadedmessagesSchema.FindOneAndUpdate(
		ctx,
		bson.M{"mainUserid": msg.Recever, "otherUserid": msg.Sender},
		bson.M{"$inc": bson.M{"numOfUnreadedMessages": 1}, "$set": bson.M{"isReaded": false}},
	).Decode(&existingRecord)

	if err != nil {
		// Record doesn't exist, create new one
		_, err = unreadedmessagesSchema.InsertOne(
			ctx,
			bson.M{"mainUserid": msg.Recever, "otherUserid": msg.Sender, "numOfUnreadedMessages": 1, "isReaded": false},
		)
		if err != nil {
			return fmt.Errorf("failed to update unread messages: %v", err)
		}
	}

	log.Printf("Message saved successfully from %s to %s", msg.Sender, msg.Recever)
	return nil
}

// helper func to check if users are friends
func (cm *ConnectionManager) isFriend(userID, friendID string) bool {
	friends := <-cm.getUserFriends(userID)
	if friends == nil {
		return false
	}
	for _, f := range friends {
		if f == friendID {
			return true
		}
	}
	return false
}
