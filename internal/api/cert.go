package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/domain"
	"github.com/zmoony/pi-gateway/internal/service"
)

type CertificateHandler struct {
	Service service.CertificateService
}

func (h CertificateHandler) List(c *gin.Context) {
	items, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if items == nil {
		items = []domain.CertificateConfig{}
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h CertificateHandler) Create(c *gin.Context) {
	var request service.UpsertCertificateConfigInput
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

func (h CertificateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid certificate id"})
		return
	}
	var request service.UpsertCertificateConfigInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrCertificateConfigNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error(), "item": item})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h CertificateHandler) Renew(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid certificate id"})
		return
	}
	item, err := h.Service.Renew(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrCertificateConfigNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error(), "item": item})
		return
	}
	c.JSON(http.StatusOK, item)
}
