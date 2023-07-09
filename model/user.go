package model

// User represents a user
type User struct {
	ID       int    `json: "id"`
	FullName string `json:"fullName"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
}
