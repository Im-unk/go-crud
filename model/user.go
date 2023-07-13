package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullName"`
	UserName string             `json:"userName"`
	Email    string             `json:"email"`
}
