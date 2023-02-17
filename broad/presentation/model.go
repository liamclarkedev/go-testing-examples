package presentation

import "github.com/google/uuid"

// UserResponse is the presentation
// representation of a User response.
type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

// UserRequest is a presentation representation of a User.
type UserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email"  binding:"required"`
}
