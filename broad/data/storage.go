package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/liamclarkedev/go-testing-examples/broad/domain"
)

var (
	// ErrUnableToBuildQuery is the annotated error when a query fails
	// to build into an SQL string.
	ErrUnableToBuildQuery = errors.New("unable to build query")

	// ErrUnableToExecuteQuery is the annotated error when a query fails
	// to execute.
	ErrUnableToExecuteQuery = errors.New("unable to execute query")
)

// Storage is the Postgres storage implementation.
type Storage struct {
	DB *sql.DB
}

// New initialises a new Storage.
func New(db *sql.DB) Storage {
	return Storage{
		DB: db,
	}
}

// Insert inserts a new User into the users table.
func (s Storage) Insert(ctx context.Context, user domain.User) (uuid.UUID, error) {
	row := UserRow(user)

	query, args, err := squirrel.
		Insert("users").
		Columns("id", "name", "email").
		Values(row.ID, row.Name, row.Email).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", domain.ErrUnknown, ErrUnableToBuildQuery)
	}

	if _, err = s.DB.ExecContext(ctx, query, args...); err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", domain.ErrAlreadyExists, ErrUnableToExecuteQuery)
	}

	return user.ID, nil
}
