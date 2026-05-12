package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/service"
)

type SystemHandler struct {
	Service service.SystemStatsService
}

func (h SystemHandler) Get(c *gin.Context) {
	stats, err := h.Service.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
