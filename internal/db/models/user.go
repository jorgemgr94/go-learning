package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserResponse represents the response after creating a user
type CreateUserResponse struct {
	ID string `json:"id"`
}

// GetUserRequest represents the request to get a user
type GetUserRequest struct {
	ID string `json:"id"`
}

// GetUserResponse represents the response when getting a user
type GetUserResponse struct {
	User *User `json:"user"`
}

// ListUsersRequest represents the request to list users
type ListUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ListUsersResponse represents the response when listing users
type ListUsersResponse struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}
