package fuzz_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/google/uuid"

	"github.com/liamclarkedev/go-testing-examples/fuzz"
)

func FuzzController_Create(f *testing.F) {
	storage := mockStorage{
		GivenID: uuid.New(),
	}

	c := fuzz.New(storage, validator.New())

	f.Add("foo", "foo@bar.com")

	// expect all errors to be handled and annotated with a known exported error.
	f.Fuzz(func(t *testing.T, name string, email string) {
		user := fuzz.User{
			Name:  name,
			Email: email,
		}
		_, err := c.Create(context.Background(), user)
		if err != nil {

			if errors.Is(err, fuzz.ErrInvalidUser) || errors.Is(err, fuzz.ErrAlreadyExists) {
				return
			}

			t.Error(err)
		}
	})
}

type mockStorage struct {
	GivenID    uuid.UUID
	GivenError error
}

func (m mockStorage) Insert(_ context.Context, _ fuzz.User) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}
