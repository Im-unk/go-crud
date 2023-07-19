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

// AddUser adds a new user
func (r *UserRepository) AddUser(user model.User) (model.User, error) {
	// First, add the user to the database
	newUser, err := r.db.AddUser(user)
	if err != nil {
		return newUser, err
	}

	// Next, index the new user data in ElasticSearch
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	err = r.searchEngine.IndexDocument(indexName, newUser)
	if err != nil {
		// If indexing fails, you may choose to handle this error accordingly,
		// like rolling back the user creation in the database.
		// For simplicity, we're not handling the error here.
		log.Println("Failed to index user in ElasticSearch:", err)
	} else {
		log.Println("User indexed successfully in ElasticSearch")
	}

	return newUser, nil
}

// UpdateUser updates a user
func (r *UserRepository) UpdateUser(user model.User) error {

	// First, update the user in the database
	r.db.UpdateUser(user)

	// Next, update the user data in ElasticSearch
	// indexName := "users" // The name of the Elasticsearch index where user data is stored.
	// err := r.searchEngine.IndexDocumentUpdate(indexName, user.ID.Hex(), user)
	// if err != nil {
	// 	// If indexing fails, you may choose to handle this error accordingly.
	// 	// For simplicity, we're not handling the error here.
	// 	log.Println("Failed to update user data in ElasticSearch:", err)
	// } else {
	// 	log.Println("User data updated successfully in ElasticSearch")
	// }

	return nil
}

// PatchUser partially updates a user
func (r *UserRepository) PatchUser(user model.User) (model.User, error) {
	return r.db.PatchUser(user)
}

// DeleteUser deletes a user by ID
func (r *UserRepository) DeleteUser(id string) error {

	// First, delete the user from the database
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid object ID format: %v", err)
	}

	newUser, err2 := r.db.GetUserByID(objID)
	if err2 != nil {
		return fmt.Errorf("Couldn't find the user ", err2)
	}

	fmt.Print(newUser.Email)

	err = r.db.DeleteUser(objID)
	if err != nil {
		return err
	}

	// Next, remove the user data from ElasticSearch
	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	err = r.searchEngine.DeleteDocumentByUniqueID(indexName, "email", newUser.Email)
	if err != nil {
		// If removing from ElasticSearch fails, you may choose to handle this error accordingly.
		// For simplicity, we're not handling the error here.
		log.Println("Failed to remove user data from ElasticSearch:", err)
	} else {
		log.Println("User data removed successfully from ElasticSearch")
	}

	return nil
}

// SearchUsers performs a search query on the user data and returns the results.
func (r *UserRepository) SearchUsers(query string) ([]search.SearchResult, error) {
	// Here, you'll use the search engine to perform the search query
	// and return the search results.
	// The implementation will vary based on your specific requirements and how you've set up the search engine.
	// In this example, we'll simply call the Search method of the search engine.

	indexName := "users" // The name of the Elasticsearch index where user data is stored.
	searchResults, err := r.searchEngine.Search(indexName, query)
	if err != nil {
		log.Println("Error searching for users in ElasticSearch:", err)
		return nil, err
	}

	log.Printf("Found %d search results for query: %s\n", len(searchResults), query)
	return searchResults, nil
}
