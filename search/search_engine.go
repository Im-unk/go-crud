package search

// SearchResult represents a single result from the search engine.
type SearchResult struct {
	ID    string  // Unique identifier of the document
	Score float64 // Relevance score of the document
	// Add other fields as needed based on your search requirements
}

// SearchEngine is an interface that defines the methods to interact with the search engine.
type SearchEngine interface {
	// IndexDocument indexes a document in the search engine.
	IndexDocument(index string, id string, data interface{}) error

	// IndexDocumentUpdate indexes a document in the search engine.
	// IndexDocumentUpdate(index, docID string, data interface{}) error

	// DeleteDocument removes a document from the search engine by its ID.
	DeleteDocument(index string, docID string) error

	// Search performs a search query on the search engine and returns the results.
	// The `query` parameter can be a string representing the search query or a more complex data structure
	// representing the query depending on your specific search requirements.
	Search(index string, query string) ([]SearchResult, error)

	// Add more methods as needed based on your search engine requirements.
}
