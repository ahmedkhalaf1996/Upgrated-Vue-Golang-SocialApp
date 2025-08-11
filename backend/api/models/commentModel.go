package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"postId" bson:"postId"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Value     string             `json:"value" bson:"value"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

type CreateComment struct {
	Value string `json:"value" bson:"value" validate:"required"`
}
