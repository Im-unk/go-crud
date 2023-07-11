package database

import "main.go/model"

// PostDatabase provides an abstraction for post-related database operations
type PostDatabase interface {
	GetPosts() ([]model.Post, error)
	GetPostByID(id int) (model.Post, error)
	AddPost(post model.Post) (model.Post, error)
	UpdatePost(post model.Post) (model.Post, error)
	PatchPost(post model.Post) (model.Post, error)
	DeletePost(id int) error
}
