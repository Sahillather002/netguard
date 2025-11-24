package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Report represents a generated report
type Report struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Format      string                 `json:"format"`
	Status      string                 `json:"status"`
	Data        map[string]interface{} `json:"data"`
	GeneratedAt time.Time              `json:"generated_at"`
	GeneratedBy string                 `json:"generated_by"`
}

// generateSecurityReport generates a security report
func generateSecurityReport(c *gin.Context) {
	reportType := c.DefaultQuery("type", "summary")
	period := c.DefaultQuery("period", "7d")

	userID, _ := c.Get("user_id")

	dataMux.RLock()
	alertCount := len(alerts)
	threatCount := len(threats)
	ruleCount := len(firewallRules)
	dataMux.RUnlock()

	// Calculate threat levels
	criticalThreats := 0
	highThreats := 0
	for _, threat := range threats {
		if threat.Severity == "critical" {
			criticalThreats++
		} else if threat.Severity == "high" {
			highThreats++
		}
	}

	report := Report{
		ID:     "RPT_" + time.Now().Format("20060102150405"),
		Name:   "Security Summary Report",
		Type:   reportType,
		Format: "json",
		Status: "completed",
		Data: map[string]interface{}{
			"period": period,
			"summary": gin.H{
				"total_alerts":      alertCount,
				"total_threats":     threatCount,
				"critical_threats":  criticalThreats,
				"high_threats":      highThreats,
				"firewall_rules":    ruleCount,
				"blocked_attempts":  156,
			},
			"top_threats": []gin.H{
				{"type": "Malware", "count": 45},
				{"type": "Brute Force", "count": 32},
				{"type": "DDoS", "count": 28},
			},
			"top_sources": []gin.H{
				{"ip": "192.168.1.100", "count": 67},
				{"ip": "10.0.0.50", "count": 45},
				{"ip": "172.16.0.25", "count": 34},
			},
			"recommendations": []string{
				"Update firewall rules to block suspicious IPs",
				"Enable additional monitoring on high-risk ports",
				"Review and update security policies",
			},
		},
		GeneratedAt: time.Now(),
		GeneratedBy: fmt.Sprintf("%v", userID),
	}

	c.JSON(http.StatusOK, report)
}

// generateThreatReport generates detailed threat report
func generateThreatReport(c *gin.Context) {
	userID, _ := c.Get("user_id")

	dataMux.RLock()
	defer dataMux.RUnlock()

	// Analyze threats
	threatsByType := make(map[string]int)
	threatsBySeverity := make(map[string]int)
	threatsByStatus := make(map[string]int)

	for _, threat := range threats {
		threatsByType[threat.Type]++
		threatsBySeverity[threat.Severity]++
		threatsByStatus[threat.Status]++
	}

	report := Report{
		ID:     "THR_RPT_" + time.Now().Format("20060102150405"),
		Name:   "Threat Analysis Report",
		Type:   "threat_analysis",
		Format: "json",
		Status: "completed",
		Data: map[string]interface{}{
			"total_threats": len(threats),
			"by_type":       threatsByType,
			"by_severity":   threatsBySeverity,
			"by_status":     threatsByStatus,
			"timeline": []gin.H{
				{"date": time.Now().AddDate(0, 0, -6).Format("2006-01-02"), "count": 12},
				{"date": time.Now().AddDate(0, 0, -5).Format("2006-01-02"), "count": 15},
				{"date": time.Now().AddDate(0, 0, -4).Format("2006-01-02"), "count": 18},
				{"date": time.Now().AddDate(0, 0, -3).Format("2006-01-02"), "count": 14},
				{"date": time.Now().AddDate(0, 0, -2).Format("2006-01-02"), "count": 20},
				{"date": time.Now().AddDate(0, 0, -1).Format("2006-01-02"), "count": 16},
				{"date": time.Now().Format("2006-01-02"), "count": len(threats)},
			},
		},
		GeneratedAt: time.Now(),
		GeneratedBy: fmt.Sprintf("%v", userID),
	}

	c.JSON(http.StatusOK, report)
}

// generateNetworkReport generates network activity report
func generateNetworkReport(c *gin.Context) {
	userID, _ := c.Get("user_id")

	report := Report{
		ID:     "NET_RPT_" + time.Now().Format("20060102150405"),
		Name:   "Network Activity Report",
		Type:   "network_activity",
		Format: "json",
		Status: "completed",
		Data: map[string]interface{}{
			"interfaces": []gin.H{
				{"name": "eth0", "packets": 1523400, "bytes": 1024000000, "errors": 0},
				{"name": "eth1", "packets": 856200, "bytes": 512000000, "errors": 2},
			},
			"top_protocols": []gin.H{
				{"protocol": "HTTPS", "percentage": 65.5},
				{"protocol": "HTTP", "percentage": 20.3},
				{"protocol": "SSH", "percentage": 8.2},
				{"protocol": "DNS", "percentage": 6.0},
			},
			"bandwidth_usage": []gin.H{
				{"hour": "00:00", "inbound": 125.5, "outbound": 87.3},
				{"hour": "06:00", "inbound": 234.2, "outbound": 156.8},
				{"hour": "12:00", "inbound": 456.7, "outbound": 298.4},
				{"hour": "18:00", "inbound": 389.1, "outbound": 245.6},
			},
		},
		GeneratedAt: time.Now(),
		GeneratedBy: fmt.Sprintf("%v", userID),
	}

	c.JSON(http.StatusOK, report)
}

// scheduleReport schedules a report generation
func scheduleReport(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Schedule  string `json:"schedule" binding:"required"`
		Format    string `json:"format"`
		Recipients []string `json:"recipients"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule := gin.H{
		"id":         "SCH_" + time.Now().Format("20060102150405"),
		"name":       req.Name,
		"type":       req.Type,
		"schedule":   req.Schedule,
		"format":     req.Format,
		"recipients": req.Recipients,
		"enabled":    true,
		"created_by": userID,
		"created_at": time.Now(),
		"next_run":   time.Now().Add(24 * time.Hour),
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Report scheduled successfully",
		"schedule": schedule,
	})
}
