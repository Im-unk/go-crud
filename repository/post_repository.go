package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (r *PostRepository) GetPostByID(id string) (model.Post, error) {

	// Convert the string ID to an ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Post{}, fmt.Errorf("invalid object ID format: %v", err)
	}

	return r.db.GetPostByID(objID)
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
func (r *PostRepository) DeletePost(id string) error {
	// Convert the string ID to an ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}
	return r.db.DeletePost(objID)
}
