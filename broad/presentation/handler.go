package presentation

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liamclarkedev/go-testing-examples/broad/domain"
)

var (
	// ErrInvalidRequest is the exported error when an invalid User is requested.
	ErrInvalidRequest = errors.New("unknown error occurred")
)

// ControllerProvider provides the domain methods.
type ControllerProvider interface {
	Create(ctx context.Context, user domain.User) (uuid.UUID, error)
}

// Handler is the presentation layer
// handler for incoming HTTP requests.
type Handler struct {
	Controller ControllerProvider
}

// New initializes a new Handler.
func New(controller ControllerProvider) Handler {
	return Handler{
		Controller: controller,
	}
}

// Create passes a new User to the storage layer.
func (h Handler) Create(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrInvalidRequest)
		return
	}

	user := domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	id, err := h.Controller.Create(c, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	resp := UserResponse{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
	}

	c.JSONP(http.StatusCreated, resp)
}
