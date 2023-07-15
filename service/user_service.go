package service

import (
	"fmt"
	"time"

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

	return s.userRepository.AddUser(user)
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
	return s.userRepository.UpdateUser(user)
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

	return s.userRepository.DeleteUser(id)
}
