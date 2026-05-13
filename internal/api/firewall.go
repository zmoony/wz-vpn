package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/zmoony/pi-gateway/internal/service"
)

type FirewallHandler struct {
	Service *service.FirewallService
	ConntrackService service.ConntrackService
}

func (h FirewallHandler) List(c *gin.Context) {
	items, err := h.Service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h FirewallHandler) Create(c *gin.Context) {
	var request service.UpsertFirewallRuleInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.Service.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h FirewallHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid firewall rule id"})
		return
	}
	var request service.UpsertFirewallRuleInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	item, err := h.Service.Update(c.Request.Context(), id, request)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrFirewallRuleNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h FirewallHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid firewall rule id"})
		return
	}
	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrFirewallRuleNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h FirewallHandler) Preview(c *gin.Context) {
	preview, err := h.Service.Preview(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"preview": preview})
}

func (h FirewallHandler) Apply(c *gin.Context) {
	state, preview, err := h.Service.Apply(c.Request.Context())
	if err != nil {
		status := http.StatusInternalServerError
		if err == service.ErrFirewallApplyPending {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pendingState": state, "preview": preview})
}

func (h FirewallHandler) Confirm(c *gin.Context) {
	if err := h.Service.Confirm(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "confirmed"})
}

func (h FirewallHandler) Pending(c *gin.Context) {
	state, err := h.Service.PendingState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pendingState": state})
}

func (h FirewallHandler) GetForwardConfig(c *gin.Context) {
	config, err := h.Service.GetForwardConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h FirewallHandler) UpdateForwardConfig(c *gin.Context) {
	var request service.UpsertFirewallForwardInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	config, err := h.Service.UpdateForwardConfig(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h FirewallHandler) Conntrack(c *gin.Context) {
	sourceIP := c.Query("source_ip")
	items, err := h.ConntrackService.List(c.Request.Context(), sourceIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
