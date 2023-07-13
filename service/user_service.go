package service

import (
	"main.go/model"
	"main.go/repository"
)

// UserService handles the user-related operations
type UserService struct {
	userRepository *repository.UserRepository
	messaging      *MessagingService // Add the messaging service as a field
}

// NewUserService creates a new UserService
func NewUserService(userRepository *repository.UserRepository, messaging *MessagingService) *UserService {
	return &UserService{
		userRepository: userRepository,
		messaging:      messaging,
	}
}

// GetUsers returns all users
func (s *UserService) GetUsers() ([]model.User, error) {
	return s.userRepository.GetUsers()
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id string) (model.User, error) {
	return s.userRepository.GetUserByID(id)
}

// AddUser adds a new user
func (s *UserService) AddUser(user model.User) (model.User, error) {
	return s.userRepository.AddUser(user)
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id string, user model.User) (model.User, error) {
	return s.userRepository.UpdateUser(user)
}

// PatchUser partially updates a user
func (s *UserService) PatchUser(id string, user model.User) (model.User, error) {
	return s.userRepository.PatchUser(user)
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id string) error {
	return s.userRepository.DeleteUser(id)
}
