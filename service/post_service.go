package service

import (
	"main.go/model"
	"main.go/repository"
)

// PostService handles the business logic for posts
type PostService struct {
	postRepository *repository.PostRepository
}

// NewPostService creates a new PostService
func NewPostService(postRepository *repository.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

// GetPosts returns all posts
func (s *PostService) GetPosts() ([]model.Post, error) {
	return s.postRepository.GetPosts()
}

// GetPostByID returns a post by ID
func (s *PostService) GetPostByID(id int) (model.Post, error) {
	return s.postRepository.GetPostByID(id)
}

// AddPost adds a new post
func (s *PostService) AddPost(post model.Post) (model.Post, error) {
	return s.postRepository.AddPost(post)
}

// UpdatePost updates a post
func (s *PostService) UpdatePost(post model.Post) (model.Post, error) {
	return s.postRepository.UpdatePost(post)
}

// PatchPost partially updates a post
func (s *PostService) PatchPost(post model.Post) (model.Post, error) {
	return s.postRepository.PatchPost(post)
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(id int) error {
	return s.postRepository.DeletePost(id)
}
