package database

import "main.go/model"

// Database provides an abstraction for the database operations
type Database interface {
	// Define the methods required for interacting with the database
	// Database defines the interface for database operations
	GetPosts() ([]model.Post, error)
	GetPostByID(id int) (model.Post, error)
	AddPost(post model.Post) (model.Post, error)
	UpdatePost(post model.Post) (model.Post, error)
	PatchPost(post model.Post) (model.Post, error)
	DeletePost(id int) error
}
