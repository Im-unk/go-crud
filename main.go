// main.go

package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/api"
	"main.go/database/implementations/mongodb"
	"main.go/repository"
	"main.go/service"
)

func main() {
	// Create a MongoDB connection
	mongoDB, err := connectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// Create the PostRepository using the MongoDB database instance
	postRepository := repository.NewPostRepository(mongodb.NewPostMongoDB(mongoDB))

	// Create the UserRepository using the MongoDB database instance
	userRepository := repository.NewUserRepository(mongodb.NewUserMongoDB(mongoDB))

	// Create the services
	postService := service.NewPostService(postRepository)
	userService := service.NewUserService(userRepository)

	// Create the API router
	router := api.NewRouter(postService, userService)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}

func connectToMongoDB() (*mongo.Database, error) {
	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping MongoDB to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Get the MongoDB database instance
	db := client.Database("project")

	// Return the MongoDB database instance
	return db, nil
}
