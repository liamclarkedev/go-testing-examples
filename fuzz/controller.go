package fuzz

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/liamclarkedev/go-testing-examples/narrow"
)

var (
	// ErrUnknown is the exported error when an unknown error occurred.
	ErrUnknown = errors.New("unknown error occurred")

	// ErrAlreadyExists is the exported error when a user already exists.
	ErrAlreadyExists = errors.New("user already exists")

	// ErrInvalidUser is the exported error when a user fails the validation check.
	ErrInvalidUser = errors.New("user invalid")
)

// StorageProvider provides the storage methods.
type StorageProvider interface {
	Insert(ctx context.Context, user User) (uuid.UUID, error)
}

// Controller is the communicator between
// the storage layer and domain.
type Controller struct {
	Storage   StorageProvider
	Validator *validator.Validate
}

// New initializes a new Controller.
func New(storage StorageProvider, validator *validator.Validate) Controller {
	return Controller{
		Storage:   storage,
		Validator: validator,
	}
}

// Create passes a new User to the storage layer.
func (c Controller) Create(ctx context.Context, user User) (uuid.UUID, error) {
	user.ID = uuid.New()

	if err := c.Validator.Struct(&user); err != nil {
		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			for _, err := range validationErr {
				return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidUser, err.Error())
			}
		}
	}

	id, err := c.Storage.Insert(ctx, user)
	if err != nil {
		if errors.Is(err, narrow.ErrUnableToExecuteQuery) {
			return uuid.Nil, ErrAlreadyExists
		}

		return uuid.Nil, ErrUnknown
	}

	return id, nil
}
