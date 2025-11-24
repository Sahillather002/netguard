package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// searchAll searches across all resources
func searchAll(c *gin.Context) {
	query := strings.ToLower(c.Query("q"))
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	results := gin.H{
		"alerts":   searchAlerts(query),
		"threats":  searchThreats(query),
		"rules":    searchFirewallRules(query),
		"users":    searchUsers(query),
		"query":    query,
	}

	c.JSON(http.StatusOK, results)
}

func searchAlerts(query string) []gin.H {
	dataMux.RLock()
	defer dataMux.RUnlock()

	var results []gin.H
	for _, alert := range alerts {
		if strings.Contains(strings.ToLower(alert.Title), query) ||
			strings.Contains(strings.ToLower(alert.Description), query) ||
			strings.Contains(strings.ToLower(alert.Severity), query) ||
			strings.Contains(strings.ToLower(alert.Source), query) {
			results = append(results, gin.H{
				"id":          alert.ID,
				"title":       alert.Title,
				"severity":    alert.Severity,
				"type":        "alert",
			})
		}
	}
	return results
}

func searchThreats(query string) []gin.H {
	dataMux.RLock()
	defer dataMux.RUnlock()

	var results []gin.H
	for _, threat := range threats {
		if strings.Contains(strings.ToLower(threat.Name), query) ||
			strings.Contains(strings.ToLower(threat.Type), query) ||
			strings.Contains(strings.ToLower(threat.SourceIP), query) ||
			strings.Contains(strings.ToLower(threat.Severity), query) {
			results = append(results, gin.H{
				"id":        threat.ID,
				"name":      threat.Name,
				"severity":  threat.Severity,
				"type":      "threat",
			})
		}
	}
	return results
}

func searchFirewallRules(query string) []gin.H {
	dataMux.RLock()
	defer dataMux.RUnlock()

	var results []gin.H
	for _, rule := range firewallRules {
		if strings.Contains(strings.ToLower(rule.Name), query) ||
			strings.Contains(strings.ToLower(rule.SourceIP), query) ||
			strings.Contains(strings.ToLower(rule.DestIP), query) ||
			strings.Contains(strings.ToLower(rule.Action), query) {
			results = append(results, gin.H{
				"id":     rule.ID,
				"name":   rule.Name,
				"action": rule.Action,
				"type":   "firewall_rule",
			})
		}
	}
	return results
}

func searchUsers(query string) []gin.H {
	usersMux.RLock()
	defer usersMux.RUnlock()

	var results []gin.H
	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Name), query) ||
			strings.Contains(strings.ToLower(user.Email), query) ||
			strings.Contains(strings.ToLower(user.Company), query) {
			results = append(results, gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"type":  "user",
			})
		}
	}
	return results
}
