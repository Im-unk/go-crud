package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"main.go/api"
	cache "main.go/cache/implementation"
	"main.go/database"
	"main.go/database/implementations/mongodb"
	"main.go/messaging"
	"main.go/repository"
	"main.go/search"
	"main.go/service"
)

func main() {
	// Create a MongoDB connection
	mongoDB, err := connectToMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	//setting up the Elastic Search Configuration
	esURL := "http://localhost:9200"
	esUsername := "root"
	esPassword := "123456"

	// Create the ElasticSearchEngine instance with the necessary configurations
	esEngine, err := search.NewElasticSearchEngine(esURL, esUsername, esPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Redis cache instance
	redisAddr := "127.0.0.1:6379"
	redisPassword := ""
	redisDB := 0

	redisCache, err := cache.NewRedisCache(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Fatal(err)
	}

	// Create the CacheService
	cacheDatabse := database.NewCacheDatabase(redisCache) // Pass the Redis cache instance

	// Create the PostRepository using the MongoDB database instance
	postRepository := repository.NewPostRepository(mongodb.NewPostMongoDB(mongoDB, cacheDatabse), esEngine)

	// Create the UserRepository using the MongoDB database instance
	userRepository := repository.NewUserRepository(mongodb.NewUserMongoDB(mongoDB, cacheDatabse), esEngine)

	// Create the NATS messaging system
	natsURL := "nats://localhost:4222" // Default URL
	natsMessaging, err := messaging.NewNatsMessaging(natsURL)
	if err != nil {
		log.Fatal(err)
	}

	// Create the MessagingService
	messagingService := service.NewMessagingService(natsMessaging)

	// Create the services
	postService := service.NewPostService(postRepository, messagingService) // Pass the messagingService
	userService := service.NewUserService(userRepository, messagingService) // Pass the messagingService

	// Create the API router
	router := api.NewRouter(postService, userService, messagingService) // Pass the messagingService

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
