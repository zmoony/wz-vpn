package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModuleStubHandler struct{}

func (h ModuleStubHandler) List(module string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"module":  module,
			"status":  "stub",
			"message": "This module is scaffolded for later implementation.",
			"items":   []any{},
		})
	}
}
