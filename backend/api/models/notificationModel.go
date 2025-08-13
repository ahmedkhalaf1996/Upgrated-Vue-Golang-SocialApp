package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Remove the User struct - we'll only store UserID now
type Notification struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Deatils   string             `json:"deatils" bson:"deatils"`
	MainUID   string             `json:"mainuid" bson:"mainuid"`
	TargetID  string             `json:"targetid" bson:"targetid"`
	UserID    string             `json:"userid" bson:"userid"` // Only store user ID
	IsReaded  bool               `json:"isreded" bson:"isreded"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
