package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/service"
)

type DDNSHandler struct {
	Service service.DDNSService
}

func (h DDNSHandler) List(c *gin.Context) {
	items, status, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if items == nil {
		items = []domain.DDNSConfig{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "runtimeStatus": status})
}

func (h DDNSHandler) Create(c *gin.Context) {
	var request service.UpsertDDNSConfigInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.Service.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "item": item})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h DDNSHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid ddns config id"})
		return
	}
	var request service.UpsertDDNSConfigInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrDDNSConfigNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error(), "item": item})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h DDNSHandler) Sync(c *gin.Context) {
	if err := h.Service.Sync(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "synced"})
}
