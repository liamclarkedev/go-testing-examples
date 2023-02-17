package data

import "github.com/google/uuid"

// UserRow is a storage representation of a User.
type UserRow struct {
	ID    uuid.UUID
	Name  string
	Email string
}
