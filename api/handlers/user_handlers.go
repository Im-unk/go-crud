package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"main.go/model"
	"main.go/service"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService service.UserService
	messaging   *service.MessagingService // Add the messaging service as a field
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService, messaging *service.MessagingService) *UserHandler {
	return &UserHandler{
		userService: userService,
		messaging:   messaging,
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
	fmt.Println("handler: Fetching user with ID:", idParam)
	user, err := h.userService.GetUserByID(idParam)
	if err != nil {
		http.Error(w, "No data found with specified ID", http.StatusNotFound)
		return
	}

	writeResponse(w, user)
}

// UpdateUser handles the PUT /users/{id} endpoint
func (h *UserHandler) UpdateUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	// id, err := strconv.Atoi(idParam) // Remove this line

	var updatedUser model.User
	err := json.NewDecoder(req.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// updatedUser.ID = id // Remove this line
	err = h.userService.UpdateUser(idParam, updatedUser) // Update the parameter to idParam
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, updatedUser) // Use updatedUser instead of user
}

// PatchUser handles the PATCH /users/{id} endpoint
func (h *UserHandler) PatchUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	// id, err := strconv.Atoi(idParam) // Remove this line

	var patchedUser model.User
	err := json.NewDecoder(req.Body).Decode(&patchedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// patchedUser.ID = id // Remove this line
	user, err := h.userService.PatchUser(idParam, patchedUser) // Update the parameter to idParam
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, user)
}

// DeleteUser handles the DELETE /users/{id} endpoint
func (h *UserHandler) DeleteUser(w http.ResponseWriter, req *http.Request) {
	idParam := mux.Vars(req)["id"]
	// id, err := strconv.Atoi(idParam) // Remove this line

	err := h.userService.DeleteUser(idParam) // Update the parameter to idParam
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
