package service

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"main.go/model"
	"main.go/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
	messaging      *MessagingService
}

func NewUserService(userRepository *repository.UserRepository, messaging *MessagingService) *UserService {
	return &UserService{
		userRepository: userRepository,
		messaging:      messaging,
	}
}

func (s *UserService) GetUsers() ([]model.User, error) {
	return s.userRepository.GetUsers()
}

func (s *UserService) GetUserByID(id string) (model.User, error) {
	// Directly call the repository method to get user by ID
	return s.userRepository.GetUserByID(id)
}

func (s *UserService) AddUser(user model.User) (model.User, error) {
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
	// Convert the string ID to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	// Update the user in the repository
	err = s.userRepository.UpdateUser(user)
	if err != nil {
		return err
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
	// Patch the user in the repository
	patchedUser, err := s.userRepository.PatchUser(user)
	if err != nil {
		return model.User{}, err
	}

	// Publish a message indicating a user has been updated
	err = s.messaging.Publish("user.updated", []byte(patchedUser.ID.Hex()))
	if err != nil {
		// Log the error, but don't affect the response
		fmt.Printf("Failed to publish user.updated message: %v\n", err)
	}

	return patchedUser, nil
}

func (s *UserService) DeleteUser(id string) error {
	// Directly call the repository method to delete user
	err := s.userRepository.DeleteUser(id)
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
