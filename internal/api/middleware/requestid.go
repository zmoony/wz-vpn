package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", time.Now().UTC().Format("20060102150405.000000000"))
		c.Next()
	}
}
