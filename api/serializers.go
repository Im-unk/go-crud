package api

import (
	"encoding/json"

	"main.go/model"
)

// UserSerializer serializes User objects to JSON
type UserSerializer struct{}

// Serialize serializes a User object to JSON
func (s *UserSerializer) Serialize(user model.User) (string, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PostSerializer serializes Post objects to JSON
type PostSerializer struct{}

// Serialize serializes a Post object to JSON
func (s *PostSerializer) Serialize(post model.Post) (string, error) {
	data, err := json.Marshal(post)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
