package repository

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	database "main.go/database/models"
	"main.go/model"
	"main.go/search"
)

// PostRepository handles the post data access
type PostRepository struct {
	db           database.PostDatabase
	searchEngine search.SearchEngine
}

// NewPostRepository creates a new PostRepository
func NewPostRepository(db database.PostDatabase, searchEngine search.SearchEngine) *PostRepository {
	return &PostRepository{
		db:           db,
		searchEngine: searchEngine,
	}
}

// GetPosts returns all posts
func (r *PostRepository) GetPosts() ([]model.Post, error) {
	return r.db.GetPosts()
}

// GetPostByID returns a post by ID
func (r *PostRepository) GetPostByID(id primitive.ObjectID) (model.Post, error) {
	return r.db.GetPostByID(id)
}

func (r *PostRepository) GetLatestInsertedPost() (model.Post, error) {
	return r.db.GetLatestInsertedPost()
}

func (r *PostRepository) AddPost(post model.Post) (model.Post, error) {
	// First, add the post to the database
	newPost, err := r.db.AddPost(post)
	if err != nil {
		return newPost, err
	}

	// Get the latest inserted post from the database
	latestPost, err := r.GetLatestInsertedPost()
	if err != nil {
		// Handle the error if necessary
		log.Println("Failed to get the latest inserted post from the database:", err)
		return newPost, nil
	}

	// Next, index the new post data in ElasticSearch with the provided "_id"
	indexName := "posts" // The name of the Elasticsearch index where post data is stored.
	err = r.searchEngine.IndexDocument(indexName, latestPost.ID.Hex(), newPost)
	if err != nil {
		log.Println("Failed to index post in ElasticSearch:", err)
	} else {
		log.Println("Post indexed successfully in ElasticSearch")
	}

	return newPost, nil
}

func (r *PostRepository) UpdatePost(post model.Post) error {
	// First, update the post in the database
	updatedPost, err := r.db.UpdatePost(post)
	if err != nil {
		return err
	}

	// Next, update the post data in Elasticsearch
	indexName := "posts" // The name of the Elasticsearch index where post data is stored.
	err = r.searchEngine.IndexDocument(indexName, updatedPost.ID.Hex(), updatedPost)
	if err != nil {
		log.Println("Failed to update post data in ElasticSearch:", err)
	} else {
		log.Println("Post data updated successfully in ElasticSearch")
	}

	return nil
}

func (r *PostRepository) PatchPost(post model.Post) (model.Post, error) {
	// First, patch the post in the database
	patchedPost, err := r.db.PatchPost(post)
	if err != nil {
		return model.Post{}, err
	}

	// Next, update the post data in Elasticsearch
	indexName := "posts" // The name of the Elasticsearch index where post data is stored.
	err = r.searchEngine.IndexDocument(indexName, patchedPost.ID.Hex(), patchedPost)
	if err != nil {
		log.Println("Failed to update post data in ElasticSearch:", err)
	} else {
		log.Println("Post data updated successfully in ElasticSearch")
	}

	return patchedPost, nil
}

func (r *PostRepository) DeletePost(id primitive.ObjectID) error {
	// First, delete the post from the database
	err := r.db.DeletePost(id)
	if err != nil {
		return err
	}

	// Next, remove the post data from ElasticSearch
	indexName := "posts" // The name of the Elasticsearch index where post data is stored.
	err = r.searchEngine.DeleteDocument(indexName, id.Hex())
	if err != nil {
		log.Println("Failed to remove post data from ElasticSearch:", err)
	} else {
		log.Println("Post data removed successfully from ElasticSearch")
	}

	return nil
}

// SearchPosts performs a search query on the post data and returns the results.
func (r *PostRepository) SearchPosts(query string) ([]search.SearchResult, error) {
	indexName := "posts" // The name of the Elasticsearch index where post data is stored.
	searchResults, err := r.searchEngine.Search(indexName, query)
	if err != nil {
		log.Println("Error searching for posts in ElasticSearch:", err)
		return nil, err
	}

	log.Printf("Found %d search results for query: %s\n", len(searchResults), query)
	return searchResults, nil
}
