package unit

import "github.com/google/uuid"

// User is a domain representation of a User.
type User struct {
	ID    uuid.UUID `validate:"required"`
	Name  string    `validate:"required,gte=3,lte=50"`
	Email string    `validate:"required,email,lte=50"`
}
