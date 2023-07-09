package main

import (
	"log"
	"net/http"

	"main.go/api"
	"main.go/repository"
	"main.go/service"
)

func main() {
	// Create the repository
	postRepository := repository.NewPostRepository()
	userRepository := repository.NewUserRepository()

	// Create the services
	postService := service.NewPostService(postRepository)
	userService := service.NewUserService(userRepository)

	// Create the API router
	router := api.NewRouter(postService, userService)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
