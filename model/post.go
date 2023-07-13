package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Post represents a single post
type Post struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title"`
	Body  string             `json:"body"`
}
