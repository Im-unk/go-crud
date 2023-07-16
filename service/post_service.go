package service

import (
	"fmt"
	"time"

	"main.go/model"
	"main.go/repository"
)

type PostService struct {
	postRepository *repository.PostRepository
	messaging      *MessagingService
	cacheService   *CacheService
}

func NewPostService(postRepository *repository.PostRepository, messaging *MessagingService, cacheService *CacheService) *PostService {
	return &PostService{
		postRepository: postRepository,
		messaging:      messaging,
		cacheService:   cacheService,
	}
}

func (s *PostService) GetPosts() ([]model.Post, error) {
	cacheKey := "posts"
	var posts []model.Post

	// Try to get posts from the cache
	err := s.cacheService.Get(cacheKey, &posts)
	if err != nil {
		// Cache miss, retrieve the posts from the repository
		posts, err = s.postRepository.GetPosts()
		if err != nil {
			return nil, err
		}

		// Store the posts in the cache
		err = s.cacheService.Set(cacheKey, posts, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set posts in cache: %v\n", err)
		}
	}

	return posts, nil
}

func (s *PostService) GetPostByID(id string) (model.Post, error) {
	cacheKey := fmt.Sprintf("post:%d", id)
	var post model.Post

	// Try to get the post from the cache
	err := s.cacheService.Get(cacheKey, &post)
	if err != nil {
		// Cache miss, retrieve the post from the repository
		post, err = s.postRepository.GetPostByID(id)
		if err != nil {
			return model.Post{}, err
		}

		// Store the post in the cache
		err = s.cacheService.Set(cacheKey, post, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set post in cache: %v\n", err)
		}
	}

	return post, nil
}

func (s *PostService) AddPost(post model.Post) (model.Post, error) {
	// Clear the posts cache
	err := s.cacheService.Delete("posts")
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete posts cache: %v\n", err)
	}

	return s.postRepository.AddPost(post)
}

func (s *PostService) UpdatePost(id int, post model.Post) (model.Post, error) {
	cacheKey := fmt.Sprintf("post:%d", id)

	// Clear the post cache
	err := s.cacheService.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete post cache: %v\n", err)
	}

	return s.postRepository.UpdatePost(post)
}

func (s *PostService) PatchPost(id int, post model.Post) (model.Post, error) {
	cacheKey := fmt.Sprintf("post:%d", id)

	// Clear the post cache
	err := s.cacheService.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete post cache: %v\n", err)
	}

	return s.postRepository.PatchPost(post)
}

func (s *PostService) DeletePost(id string) error {
	cacheKey := fmt.Sprintf("post:%d", id)

	// Clear the post cache
	err := s.cacheService.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete post cache: %v\n", err)
	}

	return s.postRepository.DeletePost(id)
}
