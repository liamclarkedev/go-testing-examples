package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/clarke94/go-testing-examples/narrow"
	"github.com/clarke94/go-testing-examples/unit"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestController_Create(t *testing.T) {
	tests := []struct {
		name    string
		storage unit.StorageProvider
		user    unit.User
		want    uuid.UUID
		wantErr error
	}{
		{
			name: "expect success given a valid User",
			storage: mockStorage{
				GivenID: uuid.Must(uuid.Parse("a6acab82-2b2e-484c-8e2b-f3f7736a26ed")),
			},
			user: unit.User{
				Name:  "foo",
				Email: "foo@bar.com",
			},
			want:    uuid.Must(uuid.Parse("a6acab82-2b2e-484c-8e2b-f3f7736a26ed")),
			wantErr: nil,
		},
		{
			name: "expect validation error given name less than 3",
			storage: mockStorage{
				GivenError: narrow.ErrUnableToExecuteQuery,
			},
			user: unit.User{
				Name:  "fo",
				Email: "foo@bar.com",
			},
			want:    uuid.Nil,
			wantErr: unit.ErrInvalidUser,
		},
		{
			name: "expect validation error given name greater than 50",
			storage: mockStorage{
				GivenError: narrow.ErrUnableToExecuteQuery,
			},
			user: unit.User{
				Name:  "this is a very long name that should be caught by the validator",
				Email: "foo@bar.com",
			},
			want:    uuid.Nil,
			wantErr: unit.ErrInvalidUser,
		},
		{
			name: "expect validation error given invalid email",
			storage: mockStorage{
				GivenError: narrow.ErrUnableToExecuteQuery,
			},
			user: unit.User{
				Name:  "foo",
				Email: "notanemail@foo",
			},
			want:    uuid.Nil,
			wantErr: unit.ErrInvalidUser,
		},
		{
			name: "expect already exists given a storage insert error",
			storage: mockStorage{
				GivenError: narrow.ErrUnableToExecuteQuery,
			},
			user: unit.User{
				Name:  "foo",
				Email: "foo@bar.com",
			},
			want:    uuid.Nil,
			wantErr: unit.ErrAlreadyExists,
		},
		{
			name: "expect unknown error given an unhandled storage error",
			storage: mockStorage{
				GivenError: errors.New("foo"),
			},
			user: unit.User{
				Name:  "foo",
				Email: "foo@bar.com",
			},
			want:    uuid.Nil,
			wantErr: unit.ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := unit.New(tt.storage, validator.New())

			got, err := c.Create(context.Background(), tt.user)

			if !cmp.Equal(err, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(err, tt.wantErr, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}
		})
	}
}

type mockStorage struct {
	GivenID    uuid.UUID
	GivenError error
}

func (m mockStorage) Insert(_ context.Context, _ unit.User) (uuid.UUID, error) {
	return m.GivenID, m.GivenError
}
