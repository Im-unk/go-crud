package search

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/olivere/elastic/v7"
)

// ElasticSearchEngine is the ElasticSearch implementation of the SearchEngine interface.
type ElasticSearchEngine struct {
	client *elastic.Client
}

// NewElasticSearchEngine creates a new instance of ElasticSearchEngine.
// The function receives the Elasticsearch URL, username, and password to create the Elasticsearch client.
func NewElasticSearchEngine(url, username, password string) (*ElasticSearchEngine, error) {
	// Create the Elasticsearch client
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetBasicAuth(username, password),
		// Add any other Elasticsearch configurations as needed.
	)
	if err != nil {
		return nil, err
	}

	return &ElasticSearchEngine{client: client}, nil
}

// IndexDocument indexes a document in the Elasticsearch index.
// func (e *ElasticSearchEngine) IndexDocument(index string, data interface{}) error {
// 	ctx := context.Background()

// 	_, err := e.client.Index().
// 		Index(index).
// 		Id("").
// 		BodyJson(data).
// 		Do(ctx)

// 	return err
// }

// IndexDocument indexes a document in the Elasticsearch index.
func (e *ElasticSearchEngine) IndexDocument(index string, data interface{}) error {
	ctx := context.Background()

	// Get the reflect value of the data to work with its fields
	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Check if the value is a struct
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct")
	}

	// Create a map to store the data fields for indexing
	docData := make(map[string]interface{})

	// Iterate over the fields of the struct and extract the field names and values
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i).Interface()
		// Use the JSON tag as the Elasticsearch field name, if available
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			// If no JSON tag is specified, use the field name as the Elasticsearch field name
			jsonTag = field.Name
		}
		docData[jsonTag] = fieldValue
	}

	_, err := e.client.Index().
		Index(index).
		BodyJson(docData).
		Do(ctx)

	return err
}

// DeleteDocument removes a document from the Elasticsearch index by its ID.
func (e *ElasticSearchEngine) DeleteDocument(index, docID string) error {
	ctx := context.Background()

	_, err := e.client.Delete().
		Index(index).
		Id(docID).
		Do(ctx)

	return err
}

// DeleteDocumentByUniqueID removes a document from the Elasticsearch index by a unique identifier (e.g., email or username).
func (e *ElasticSearchEngine) DeleteDocumentByUniqueID(index string, uniqueIDField string, uniqueIDValue string) error {
	ctx := context.Background()

	// Implement the logic to find the document by the unique identifier and get its ID.
	// For example, if uniqueIDField is "email", you can use a search query to find the document by the email value.
	// Extract the document ID and then delete the document using the DeleteDocument function.

	// Example code (replace with your actual logic):
	query := elastic.NewTermQuery(uniqueIDField, uniqueIDValue)
	searchResult, err := e.client.Search().Index(index).Query(query).Do(ctx)
	if err != nil {
		return err
	}

	if searchResult.TotalHits() == 0 {
		return fmt.Errorf("document with unique identifier not found")
	}

	// Assuming there's only one document matching the unique identifier, get its ID.
	docID := searchResult.Hits.Hits[0].Id

	// Delete the document using the DeleteDocument function.
	_, err = e.client.Delete().
		Index(index).
		Id(docID).
		Do(ctx)

	return err
}

// Search performs a search query on the Elasticsearch index and returns the results.
func (e *ElasticSearchEngine) Search(index, query string) ([]SearchResult, error) {
	ctx := context.Background()

	// Perform the search query
	result, err := e.client.Search(index).
		Query(elastic.NewQueryStringQuery(query)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	// Extract the search results
	var searchResults []SearchResult
	for _, hit := range result.Hits.Hits {
		var sr SearchResult
		err := json.Unmarshal(hit.Source, &sr)
		if err != nil {
			return nil, err
		}
		searchResults = append(searchResults, sr)
	}

	return searchResults, nil
}

// Add more methods as needed based on your Elasticsearch requirements.
