package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/service"
)

type SettingsHandler struct {
	Service service.SettingsService
}

func (h SettingsHandler) Get(c *gin.Context) {
	view, err := h.Service.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, view)
}

func (h SettingsHandler) Update(c *gin.Context) {
	var request service.UpdateSettingsInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	view, err := h.Service.Update(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, view)
}
