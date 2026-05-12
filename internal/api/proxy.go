package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/service"
)

type ProxyHandler struct {
	Service service.ProxyService
}

func (h ProxyHandler) List(c *gin.Context) {
	items, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h ProxyHandler) Create(c *gin.Context) {
	var request service.UpsertProxyHostInput
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

func (h ProxyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid proxy id"})
		return
	}
	var request service.UpsertProxyHostInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrProxyHostNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error(), "item": item})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h ProxyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid proxy id"})
		return
	}
	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrProxyHostNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
