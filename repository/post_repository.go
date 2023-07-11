package repository

import (
	database "main.go/database/models"
	"main.go/model"
)

// PostRepository handles the post data access
type PostRepository struct {
	db database.PostDatabase
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db database.PostDatabase) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

// GetPosts returns all posts
func (r *PostRepository) GetPosts() ([]model.Post, error) {
	return r.db.GetPosts()
}

// GetPostByID returns a post by ID
func (r *PostRepository) GetPostByID(id int) (model.Post, error) {
	return r.db.GetPostByID(id)
}

// AddPost adds a new post
func (r *PostRepository) AddPost(post model.Post) (model.Post, error) {
	return r.db.AddPost(post)
}

// UpdatePost updates a post
func (r *PostRepository) UpdatePost(post model.Post) (model.Post, error) {
	return r.db.UpdatePost(post)
}

// PatchPost partially updates a post
func (r *PostRepository) PatchPost(post model.Post) (model.Post, error) {
	return r.db.PatchPost(post)
}

// DeletePost deletes a post by ID
func (r *PostRepository) DeletePost(id int) error {
	return r.db.DeletePost(id)
}
