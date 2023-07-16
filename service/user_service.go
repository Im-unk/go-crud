package service

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/model"
	"main.go/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
	messaging      *MessagingService
	cacheService   *CacheService
}

func NewUserService(userRepository *repository.UserRepository, messaging *MessagingService, cacheService *CacheService) *UserService {
	return &UserService{
		userRepository: userRepository,
		messaging:      messaging,
		cacheService:   cacheService,
	}
}

func (s *UserService) GetUsers() ([]model.User, error) {
	cacheKey := "users"
	var users []model.User

	// Try to get users from the cache
	err := s.cacheService.Get(cacheKey, &users)
	if err != nil {
		// Cache miss, retrieve the users from the repository
		users, err = s.userRepository.GetUsers()
		if err != nil {
			return nil, err
		}

		// Store the users in the cache
		err = s.cacheService.Set(cacheKey, users, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set users in cache: %v\n", err)
		}
	}

	return users, nil
}

func (s *UserService) GetUserByID(id string) (model.User, error) {
	cacheKey := fmt.Sprintf("user:%s", id)
	var user model.User
	fmt.Println("service: Fetching user with ID:", id)

	// Try to get the user from the cache
	err := s.cacheService.Get(cacheKey, &user)
	if err != nil {
		// Cache miss, retrieve the user from the repository
		user, err = s.userRepository.GetUserByID(id)
		if err != nil {
			return model.User{}, err
		}

		// Store the user in the cache
		err = s.cacheService.Set(cacheKey, user, time.Hour)
		if err != nil {
			// Log the error, but don't affect the response
			fmt.Printf("Failed to set user in cache: %v\n", err)
		}
	}

	return user, nil
}

func (s *UserService) AddUser(user model.User) (model.User, error) {
	// Clear the users cache
	err := s.cacheService.Delete("users")
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete users cache: %v\n", err)
	}

	// Add the user to the repository
	addedUser, err := s.userRepository.AddUser(user)
	if err != nil {
		return model.User{}, err
	}

	// Convert the primitive.ObjectID to a string
	userID := addedUser.ID.Hex()

	// Publish a message indicating a new user has been added
	err = s.messaging.Publish("user.added", []byte(userID))
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to publish user.added message: %v\n", err)
	}

	return addedUser, nil
}

func (s *UserService) UpdateUser(id string, user model.User) error {
	cacheKey := fmt.Sprintf("user:%s", id)

	// Clear the user cache
	err := s.cacheService.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete user cache: %v\n", err)
	}

	// Update the user in the repository
	err = s.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}

	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	// Publish a message indicating a user has been updated
	err = s.messaging.Publish("user.updated", []byte(objID.Hex()))
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to publish user.updated message: %v\n", err)
	}

	return nil
}

func (s *UserService) PatchUser(id string, user model.User) (model.User, error) {
	cacheKey := fmt.Sprintf("user:%s", id)

	// Clear the user cache
	err := s.cacheService.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete user cache: %v\n", err)
	}

	return s.userRepository.PatchUser(user)
}

func (s *UserService) DeleteUser(id string) error {
	cacheKey := fmt.Sprintf("user:%s", id)

	// Clear the user cache
	err := s.cacheService.Delete(cacheKey)
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to delete user cache: %v\n", err)
	}

	// Delete the user from the repository
	err = s.userRepository.DeleteUser(id)
	if err != nil {
		return err
	}

	// Publish a message indicating a user has been deleted
	err = s.messaging.Publish("user.deleted", []byte(id))
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to publish user.deleted message: %v\n", err)
	}

	return nil
}
