package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"main.go/model"
	"main.go/service"
)

// PostHandler handles HTTP requests for posts
type PostHandler struct {
	postService service.PostService
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
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
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	post, err := h.postService.GetPostByID(id)
	if err != nil {
		http.Error(w, "No data found with specified ID", http.StatusNotFound)
		return
	}

	writeResponse(w, post)
}

// UpdatePost handles the PUT /posts/{id} endpoint
func (h *PostHandler) UpdatePost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	var updatedPost model.Post
	err = json.NewDecoder(req.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//updatedPost.ID = id
	post, err := h.postService.UpdatePost(id, updatedPost) // Pass the ID as a separate parameter
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeResponse(w, post)
}

// PatchPost handles the PATCH /posts/{id} endpoint
func (h *PostHandler) PatchPost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	var patchedPost model.Post
	err = json.NewDecoder(req.Body).Decode(&patchedPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//patchedPost.ID = id
	post, err := h.postService.PatchPost(id, patchedPost) // Pass the ID as a separate parameter
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, post)
}

// DeletePost handles the DELETE /posts/{id} endpoint
func (h *PostHandler) DeletePost(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	err = h.postService.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
