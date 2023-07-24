package api

import (
	"net/http"

	"github.com/gorilla/mux"
	api "main.go/api/handlers"
	"main.go/service"
)

// Router handles the API routing
type Router struct {
	router      *mux.Router
	postHandler *api.PostHandler
	userHandler *api.UserHandler
}

// NewRouter creates a new API router
func NewRouter(postService *service.PostService, userService *service.UserService, messagingService *service.MessagingService) *Router {
	router := mux.NewRouter()

	postHandler := api.NewPostHandler(*postService, messagingService)
	userHandler := api.NewUserHandler(*userService, messagingService)

	// Register API endpoints
	router.HandleFunc("/posts", postHandler.GetPosts).Methods("GET")
	router.HandleFunc("/posts", postHandler.AddPost).Methods("POST")
	router.HandleFunc("/posts/{id}", postHandler.GetPost).Methods("GET")
	router.HandleFunc("/posts/{id}", postHandler.UpdatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", postHandler.PatchPost).Methods("PATCH")
	router.HandleFunc("/posts/{id}", postHandler.DeletePost).Methods("DELETE")

	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.AddUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.PatchUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	router.HandleFunc("/users/search/{query}", userHandler.SearchUser).Methods("GET")

	return &Router{
		router:      router,
		postHandler: postHandler,
		userHandler: userHandler,
	}
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
