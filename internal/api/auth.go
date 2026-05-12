package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/config"
	"github.com/zmoony/pi-gateway/internal/service"
)

type AuthHandler struct {
	Service service.AuthService
	Config  config.Config
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h AuthHandler) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, user, err := h.Service.Login(c.Request.Context(), request.Username, request.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.SetCookie("pi_gateway_token", token, h.Config.SessionTTLHours*3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
		"expiresAt": time.Now().Add(time.Duration(h.Config.SessionTTLHours) * time.Hour),
	})
}

func (h AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("pi_gateway_token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h AuthHandler) Me(c *gin.Context) {
	username, _ := c.Get("username")
	userID, _ := c.Get("userID")
	c.JSON(http.StatusOK, gin.H{
		"id":       userID,
		"username": username,
	})
}
