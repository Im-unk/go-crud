package repository

import (
	database "main.go/database/models"
	"main.go/model"
)

// UserRepository handles the user data access
type UserRepository struct {
	db database.UserDatabase
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db database.UserDatabase) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUsers returns all users
func (r *UserRepository) GetUsers() ([]model.User, error) {
	return r.db.GetUsers()
}

// GetUserByID returns a user by ID
func (r *UserRepository) GetUserByID(id string) (model.User, error) {
	return r.db.GetUserByID(id)
}

// AddUser adds a new user
func (r *UserRepository) AddUser(user model.User) (model.User, error) {
	return r.db.AddUser(user)
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(user model.User) (model.User, error) {
	return r.db.UpdateUser(user)
}

// PatchUser partially updates a user
func (r *UserRepository) PatchUser(user model.User) (model.User, error) {
	return r.db.PatchUser(user)
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) error {
	return r.db.DeleteUser(id)
}
