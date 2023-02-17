package presentation

import (
	"github.com/gin-gonic/gin"
)

// Routes defines the HTTP routes with the given gin Router.
func Routes(router *gin.Engine, handler Handler) *gin.Engine {
	router.POST("/v1/user", handler.Create)

	return router
}
