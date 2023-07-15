package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Convert the string ID to an ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, fmt.Errorf("invalid object ID format: %v", err)
	}

	return r.db.GetUserByID(objID)
}

// AddUser adds a new user
func (r *UserRepository) AddUser(user model.User) (model.User, error) {
	return r.db.AddUser(user)
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(user model.User) error {
	fmt.Println("repo: Updating user with ID:", user.ID)

	filter := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"fullname": user.FullName,
			"username": user.UserName,
			"email":    user.Email,
		},
	}

	return r.db.UpdateUser(filter, update)
}

// PatchUser partially updates a user
func (r *UserRepository) PatchUser(user model.User) (model.User, error) {
	return r.db.PatchUser(user)
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) error {
	fmt.Println("repo: Deleting user with ID:", id)

	// Convert the string ID to an ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	return r.db.DeleteUser(objID)
}
