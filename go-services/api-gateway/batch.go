package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// batchDeleteAlerts deletes multiple alerts
func batchDeleteAlerts(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataMux.Lock()
	deleted := 0
	for _, id := range req.IDs {
		if _, exists := alerts[id]; exists {
			delete(alerts, id)
			deleted++
		}
	}
	dataMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Alerts deleted successfully",
		"deleted": deleted,
		"total":   len(req.IDs),
	})
}

// batchUpdateAlerts updates multiple alerts status
func batchUpdateAlerts(c *gin.Context) {
	var req struct {
		IDs    []string `json:"ids" binding:"required"`
		Status string   `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataMux.Lock()
	updated := 0
	for _, id := range req.IDs {
		if alert, exists := alerts[id]; exists {
			alert.Status = req.Status
			updated++
		}
	}
	dataMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Alerts updated successfully",
		"updated": updated,
		"total":   len(req.IDs),
	})
}

// batchDeleteThreats deletes multiple threats
func batchDeleteThreats(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataMux.Lock()
	deleted := 0
	for _, id := range req.IDs {
		if _, exists := threats[id]; exists {
			delete(threats, id)
			deleted++
		}
	}
	dataMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Threats deleted successfully",
		"deleted": deleted,
		"total":   len(req.IDs),
	})
}

// batchDeleteFirewallRules deletes multiple firewall rules
func batchDeleteFirewallRules(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataMux.Lock()
	deleted := 0
	for _, id := range req.IDs {
		if _, exists := firewallRules[id]; exists {
			delete(firewallRules, id)
			deleted++
		}
	}
	dataMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"message": "Firewall rules deleted successfully",
		"deleted": deleted,
		"total":   len(req.IDs),
	})
}

// batchEnableFirewallRules enables/disables multiple firewall rules
func batchEnableFirewallRules(c *gin.Context) {
	var req struct {
		IDs     []string `json:"ids" binding:"required"`
		Enabled bool     `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataMux.Lock()
	updated := 0
	for _, id := range req.IDs {
		if rule, exists := firewallRules[id]; exists {
			rule.Enabled = req.Enabled
			updated++
		}
	}
	dataMux.Unlock()

	action := "enabled"
	if !req.Enabled {
		action = "disabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Firewall rules " + action + " successfully",
		"updated": updated,
		"total":   len(req.IDs),
	})
}
