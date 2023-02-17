package e2e_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/liamclarkedev/go-testing-examples/broad/data"
	"github.com/liamclarkedev/go-testing-examples/broad/domain"
	"github.com/liamclarkedev/go-testing-examples/broad/presentation"
)

var db *sql.DB

func TestE2E(t *testing.T) {
	storage := data.New(db)
	controller := domain.New(storage, validator.New())

	validRequest := presentation.UserRequest{
		Name:  "Foo Bar",
		Email: "newuser@bar.com",
	}

	validBody, err := json.Marshal(validRequest)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		router     *gin.Engine
		handler    presentation.Handler
		method     string
		url        string
		body       io.Reader
		wantStatus string
	}{
		{
			name:       "expect User given valid request",
			router:     gin.New(),
			handler:    presentation.New(controller),
			method:     http.MethodPost,
			url:        "/v1/user",
			body:       bytes.NewReader(validBody),
			wantStatus: http.StatusText(http.StatusCreated),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := presentation.Routes(tt.router, tt.handler)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			if err != nil {
				t.Error(err)
			}

			router.ServeHTTP(w, req)

			gotStatus := http.StatusText(w.Result().StatusCode)

			if !cmp.Equal(gotStatus, tt.wantStatus) {
				t.Error(cmp.Diff(gotStatus, tt.wantStatus))
			}
		})
	}
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
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
