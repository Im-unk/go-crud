package repository

import (
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	database "main.go/database/models"
	"main.go/model"
	"main.go/search"
)

// UserRepository handles the user data access
type UserRepository struct {
	db           database.UserDatabase
	searchEngine search.SearchEngine
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db database.UserDatabase, searchEngine search.SearchEngine) *UserRepository {
	return &UserRepository{
		db:           db,
		searchEngine: searchEngine,
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

func (r *UserRepository) GetLatestInsertedUser() (model.User, error) {
	return r.db.GetLatestInsertedUser()
}

func (r *UserRepository) AddUser(user model.User) (model.User, error) {
	// First, add the user to the database
	newUser, err := r.db.AddUser(user)
	if err != nil {
		return newUser, err
	}

	// Get the latest inserted user from the database
	latestUser, err := r.GetLatestInsertedUser()
	if err != nil {
		// Handle the error if necessary
		log.Println("Failed to get the latest inserted user from the database:", err)
		return newUser, nil
	}

	// Next, index the new user data in ElasticSearch with the provided "_id"
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	err = r.searchEngine.IndexDocument(indexName, latestUser.ID.Hex(), newUser)
	if err != nil {
		log.Println("Failed to index user in ElasticSearch:", err)
	} else {
		log.Println("User indexed successfully in ElasticSearch")
	}

	return newUser, nil
}

func (r *UserRepository) UpdateUser(user model.User) error {
	// First, update the user in the database
	updatedUser, err := r.db.UpdateUser(user)
	if err != nil {
		return err
	}

	// Next, update the user data in Elasticsearch
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	err = r.searchEngine.IndexDocument(indexName, updatedUser.ID.Hex(), updatedUser)
	if err != nil {
		log.Println("Failed to update user data in ElasticSearch:", err)
	} else {
		log.Println("User data updated successfully in ElasticSearch")
	}

	return nil
}

func (r *UserRepository) PatchUser(user model.User) (model.User, error) {
	// First, patch the user in the database
	patchedUser, err := r.db.PatchUser(user)
	if err != nil {
		return model.User{}, err
	}

	// Next, update the user data in Elasticsearch
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	err = r.searchEngine.IndexDocument(indexName, patchedUser.ID.Hex(), patchedUser)
	if err != nil {
		log.Println("Failed to update user data in ElasticSearch:", err)
	} else {
		log.Println("User data updated successfully in ElasticSearch")
	}

	return patchedUser, nil
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) error {
	// Convert the string ID to an ObjectId
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	// First, delete the user from the database
	err = r.db.DeleteUser(objID)
	if err != nil {
		return err
	}

	// Next, remove the user data from ElasticSearch
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	err = r.searchEngine.DeleteDocument(indexName, id)
	if err != nil {
		log.Println("Failed to remove user data from ElasticSearch:", err)
	} else {
		log.Println("User data removed successfully from ElasticSearch")
	}

	return nil
}

// SearchUsers performs a search query on the user data and returns the results.
func (r *UserRepository) SearchUsers(query string) ([]search.SearchResult, error) {
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	searchResults, err := r.searchEngine.Search(indexName, query)
	if err != nil {
		log.Println("Error searching for users in ElasticSearch:", err)
		return nil, err
	}

	log.Printf("Found %d search results for query: %s\n", len(searchResults), query)
	return searchResults, nil
}
