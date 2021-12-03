package narrow

import "github.com/google/uuid"

// User is a storage representation of a User.
type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}
