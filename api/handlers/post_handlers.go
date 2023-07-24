package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"main.go/model"
	"main.go/service"
)

// PostHandler handles HTTP requests for posts
type PostHandler struct {
	postService service.PostService
	messaging   *service.MessagingService // Add the messaging service as a field
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(postService service.PostService, messaging *service.MessagingService) *PostHandler {
	return &PostHandler{
		postService: postService,
		messaging:   messaging,
	}
}

// GetPosts handles the GET /posts endpoint
func (h *PostHandler) GetPosts(w http.ResponseWriter, req *http.Request) {
	posts, err := h.postService.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, posts)
}

// AddPost handles the POST /posts endpoint
func (h *PostHandler) AddPost(w http.ResponseWriter, req *http.Request) {
	var newPost model.Post
	err := json.NewDecoder(req.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.postService.AddPost(newPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, post)
}

// GetPost handles the GET /posts/{id} endpoint
func (h *PostHandler) GetPost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	fmt.Println("handler: Fetching post with ID:", idParam)

	post, err := h.postService.GetPostByID(idParam)
	if err != nil {
		http.Error(w, "No data found with specified ID", http.StatusNotFound)
		return
	}

	writeResponse(w, post)
}

// UpdatePost handles the PUT /posts/{id} endpoint
func (h *PostHandler) UpdatePost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]

	var updatedPost model.Post
	err := json.NewDecoder(req.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.postService.UpdatePost(idParam, updatedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, updatedPost)
}

// PatchPost handles the PATCH /posts/{id} endpoint
func (h *PostHandler) PatchPost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]

	var patchedPost model.Post
	err := json.NewDecoder(req.Body).Decode(&patchedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post, err := h.postService.PatchPost(idParam, patchedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, post)
}

// DeletePost handles the DELETE /posts/{id} endpoint
func (h *PostHandler) DeletePost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]

	err := h.postService.DeletePost(idParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SearchPost handles the GET /search/posts/{query} endpoint
func (h *PostHandler) SearchPost(w http.ResponseWriter, req *http.Request) {
	queryParam := mux.Vars(req)["query"]

	results, err := h.postService.SearchPost(queryParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send it as JSON in the response
	json.NewEncoder(w).Encode(results)
}
