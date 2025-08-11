package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostModel struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Creator      string             `json:"creator" bson:"creator"`
	Title        string             `json:"title" bson:"title"`
	Message      string             `json:"message" bson:"message"`
	Name         string             `json:"name" bson:"name"`
	SelectedFile string             `json:"selectedFile" bson:"selectedFile"`
	Likes        []string           `json:"likes" bson:"likes"`
	Comments     []CommentWithUser  `json:"comments,omitempty" bson:"comments,omitempty"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

// Add this struct for comments with user data
type CommentWithUser struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	PostID    primitive.ObjectID `json:"postId" bson:"postId"`
	UserID    primitive.ObjectID `json:"userId" bson:"userId"`
	Value     string             `json:"value" bson:"value"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	User      User               `json:"user" bson:"user"`
}

// interfaces
type CreateOrUpdatePost struct {
	Title        string `json:"title" bson:"title" validate:"required"`
	Message      string `json:"message" bson:"message" validate:"required,min=5"`
	SelectedFile string `json:"selectedFile" bson:"selectedFile"`
}

// interfaces
type ComnmentPost struct {
	Value string `json:"value" bson:"value" validate:"required"`
}
