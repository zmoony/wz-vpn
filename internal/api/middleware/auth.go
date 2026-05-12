package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/service"
)

func Auth(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("pi_gateway_token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authentication required"})
			return
		}

		claims, err := authService.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("username", claims["username"])
		c.Next()
	}
}
