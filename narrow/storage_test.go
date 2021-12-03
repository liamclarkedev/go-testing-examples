package narrow_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/clarke94/go-testing-examples/narrow"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB

func TestStorage_Insert(t *testing.T) {
	tests := []struct {
		name    string
		db      *sql.DB
		ctx     context.Context
		user    narrow.User
		want    uuid.UUID
		wantErr error
	}{
		{
			name: "expect success when inserting a new User",
			db:   db,
			ctx:  context.Background(),
			user: narrow.User{
				ID:    uuid.Must(uuid.Parse("5db0da2f-170f-4912-b43c-e2e2da009bdd")),
				Name:  "foo",
				Email: "new@bar.com",
			},
			want:    uuid.Must(uuid.Parse("5db0da2f-170f-4912-b43c-e2e2da009bdd")),
			wantErr: nil,
		},
		{
			name: "expect fail when inserting an email that already exists",
			db:   db,
			ctx:  context.Background(),
			user: narrow.User{
				ID:    uuid.New(),
				Name:  "foo",
				Email: "foo@bar.com",
			},
			want:    uuid.Nil,
			wantErr: narrow.ErrUnableToExecuteQuery,
		},
		{
			name: "expect fail when inserting an ID that already exists",
			db:   db,
			ctx:  context.Background(),
			user: narrow.User{
				ID:    uuid.Must(uuid.Parse("a6acab82-2b2e-484c-8e2b-f3f7736a26ed")),
				Name:  "foo",
				Email: "id@bar.com",
			},
			want:    uuid.Nil,
			wantErr: narrow.ErrUnableToExecuteQuery,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := narrow.New(tt.db)

			got, gotErr := s.Insert(tt.ctx, tt.user)

			if !cmp.Equal(got, tt.want) {
				t.Error(cmp.Diff(got, tt.want))
			}

			if !cmp.Equal(gotErr, tt.wantErr, cmpopts.EquateErrors()) {
				t.Error(cmp.Diff(gotErr, tt.wantErr, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	options := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=user",
			"POSTGRES_PASSWORD=secret",
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5432"},
			},
		},
		Mounts: []string{fmt.Sprintf("%s/stub:/docker-entrypoint-initdb.d", path)},
	}

	resource, err := pool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	err = resource.Expire(30)
	if err != nil {
		log.Fatalf("Could not expire resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open(
			"postgres",
			"postgres://user:secret@localhost:5432?sslmode=disable",
		)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to narrow: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
