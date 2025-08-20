package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

// CachedResponse represents the cached data structure
type CachedGetAllPostResponse struct {
	Data          []bson.M `json:"data"`
	CurrentPage   int      `json:"currentPage"`
	NumberOfPages float64  `json:"numberOfPages"`
}

// CachedGetUserResponse represents the cached data structure for user profile
type CachedGetUserResponse struct {
	User          UserModel `json:"user"`
	Posts         []bson.M  `json:"posts"`
	CurrentPage   int       `json:"currentPage"`
	NumberOfPages float64   `json:"numberOfPages"`
}
