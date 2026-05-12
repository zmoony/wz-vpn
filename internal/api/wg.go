package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/service"
)

type WireGuardHandler struct {
	Service service.WireGuardService
}

type createPeerRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h WireGuardHandler) List(c *gin.Context) {
	peers, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": peers})
}

func (h WireGuardHandler) Create(c *gin.Context) {
	var request createPeerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.Service.Create(c.Request.Context(), service.CreatePeerInput{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h WireGuardHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid peer id"})
		return
	}

	peer, err := h.Service.Toggle(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrPeerNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, peer)
}

func (h WireGuardHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid peer id"})
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrPeerNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
