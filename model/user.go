package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// User represents a user
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullname" bson:"fullname"`
	UserName string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
}
