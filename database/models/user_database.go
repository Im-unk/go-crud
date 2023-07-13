package database

import "main.go/model"

// UserDatabase provides an abstraction for user-related database operations
type UserDatabase interface {
	GetUsers() ([]model.User, error)
	GetUserByID(id string) (model.User, error)
	AddUser(user model.User) (model.User, error)
	UpdateUser(user model.User) (model.User, error)
	PatchUser(post model.User) (model.User, error)
	DeleteUser(id string) error
}
