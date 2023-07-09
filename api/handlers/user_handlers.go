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
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers handles the GET /users endpoint
func (h *UserHandler) GetUsers(w http.ResponseWriter, req *http.Request) {
	users, err := h.userService.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, users)
}

// AddUser handles the POST /users endpoint
func (h *UserHandler) AddUser(w http.ResponseWriter, req *http.Request) {
	var newUser model.User
	err := json.NewDecoder(req.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.AddUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, user)
}

// GetUser handles the GET /users/{id} endpoint
func (h *UserHandler) GetUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, "No data found with specified ID", http.StatusNotFound)
		return
	}

	writeResponse(w, user)
}

// UpdateUser handles the PUT /users/{id} endpoint
func (h *UserHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	var updatedUser model.User
	err = json.NewDecoder(req.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser.ID = id
	user, err := h.userService.UpdateUser(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, user)
}

// PatchUser handles the PATCH /users/{id} endpoint
func (h *UserHandler) PatchUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	var patchedUser model.User
	err = json.NewDecoder(req.Body).Decode(&patchedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	patchedUser.ID = id
	user, err := h.userService.PatchUser(patchedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, user)
}

// DeleteUser handles the DELETE /users/{id} endpoint
func (h *UserHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "ID could not be converted to integer", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}