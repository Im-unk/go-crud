package service

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/model"
	"main.go/repository"
	"main.go/search"
)

type PostService struct {
	postRepository *repository.PostRepository
	messaging      *MessagingService
}

func NewPostService(postRepository *repository.PostRepository, messaging *MessagingService) *PostService {
	return &PostService{
		postRepository: postRepository,
		messaging:      messaging,
	}
}

func (s *PostService) GetPosts() ([]model.Post, error) {
	return s.postRepository.GetPosts()
}

func (s *PostService) GetPostByID(id string) (model.Post, error) {
	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Post{}, fmt.Errorf("invalid object ID format: %v", err)
	}

	// Directly call the repository method to get post by ID
	return s.postRepository.GetPostByID(objID)
}

func (s *PostService) AddPost(post model.Post) (model.Post, error) {
	addedPost, err := s.postRepository.AddPost(post)
	if err != nil {
		return model.Post{}, err
	}

	// Get the latest inserted post from the database
	latestPost, err := s.postRepository.GetLatestInsertedPost()
	if err != nil {
		// Handle the error if necessary
		log.Println("Failed to get the latest inserted post from the database:", err)
		return latestPost, nil
	}

	// Convert the primitive.ObjectID to a string
	postID := latestPost.ID.Hex()

	// Publish a message indicating a new post has been added
	err = s.messaging.Publish("post.added", []byte(postID))
	if err != nil {
		// Log the error if publishing fails
		fmt.Printf("Failed to publish post.added message: %v\n", err)
	} else {
		// Log success if publishing is successful
		fmt.Printf("Successfully published post.added message for post %s\n", postID)
	}

	return addedPost, nil
}

func (s *PostService) UpdatePost(id string, post model.Post) error {
	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	// Update the post in the repository
	err = s.postRepository.UpdatePost(post)
	if err != nil {
		return err
	}

	// Publish a message indicating a post has been updated
	err = s.messaging.Publish("post.updated", []byte(objID.Hex()))
	if err != nil {
		// Log the error if publishing fails
		fmt.Printf("Failed to publish post.updated message: %v\n", err)
	} else {
		// Log success if publishing is successful
		fmt.Printf("Successfully published post.updated message for post %s\n", objID.Hex())
	}

	return nil
}

func (s *PostService) PatchPost(id string, post model.Post) (model.Post, error) {
	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Post{}, fmt.Errorf("invalid object ID format: %v", err)
	}

	// Patch the post in the repository
	patchedPost, err := s.postRepository.PatchPost(post)
	if err != nil {
		return model.Post{}, err
	}

	// Publish a message indicating a post has been updated
	err = s.messaging.Publish("post.updated", []byte(objID.Hex()))
	if err != nil {
		// Log the error if publishing fails
		fmt.Printf("Failed to publish post.updated message: %v\n", err)
	} else {
		// Log success if publishing is successful
		fmt.Printf("Successfully published post.updated message for post %s\n", objID.Hex())
	}

	return patchedPost, nil
}

func (s *PostService) DeletePost(id string) error {
	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	// Directly call the repository method to delete post
	err = s.postRepository.DeletePost(objID)
	if err != nil {
		return err
	}

	// Publish a message indicating a post has been deleted
	err = s.messaging.Publish("post.deleted", []byte(objID.Hex()))
	if err != nil {
		// Log the error if publishing fails
		fmt.Printf("Failed to publish post.deleted message: %v\n", err)
	} else {
		// Log success if publishing is successful
		fmt.Printf("Successfully published post.deleted message for post %s\n", objID.Hex())
	}

	return nil
}

func (s *PostService) SearchPost(query string) ([]search.SearchResult, error) {
	// Directly call the repository method to search for posts
	results, err := s.postRepository.SearchPosts(query)
	if err != nil {
		return nil, err
	}

	return results, nil
}
