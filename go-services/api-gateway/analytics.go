package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// getAlertAnalytics returns alert analytics
func getAlertAnalytics(c *gin.Context) {
	dataMux.RLock()
	defer dataMux.RUnlock()

	// Count by severity
	severityCounts := make(map[string]int)
	// Count by status
	statusCounts := make(map[string]int)
	// Count by source
	sourceCounts := make(map[string]int)
	// Count by hour
	hourCounts := make(map[int]int)

	for _, alert := range alerts {
		severityCounts[alert.Severity]++
		statusCounts[alert.Status]++
		sourceCounts[alert.Source]++
		hourCounts[alert.Timestamp.Hour()]++
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       len(alerts),
		"by_severity": severityCounts,
		"by_status":   statusCounts,
		"by_source":   sourceCounts,
		"by_hour":     hourCounts,
	})
}

// getThreatAnalytics returns threat analytics
func getThreatAnalytics(c *gin.Context) {
	dataMux.RLock()
	defer dataMux.RUnlock()

	// Count by type
	typeCounts := make(map[string]int)
	// Count by severity
	severityCounts := make(map[string]int)
	// Count by status
	statusCounts := make(map[string]int)
	// Top source IPs
	sourceIPCounts := make(map[string]int)
	// Total detections
	totalDetections := 0

	for _, threat := range threats {
		typeCounts[threat.Type]++
		severityCounts[threat.Severity]++
		statusCounts[threat.Status]++
		sourceIPCounts[threat.SourceIP]++
		totalDetections += threat.Detections
	}

	c.JSON(http.StatusOK, gin.H{
		"total":            len(threats),
		"total_detections": totalDetections,
		"by_type":          typeCounts,
		"by_severity":      severityCounts,
		"by_status":        statusCounts,
		"top_source_ips":   sourceIPCounts,
	})
}

// getFirewallAnalytics returns firewall analytics
func getFirewallAnalytics(c *gin.Context) {
	dataMux.RLock()
	defer dataMux.RUnlock()

	// Count by action
	actionCounts := make(map[string]int)
	// Count by protocol
	protocolCounts := make(map[string]int)
	// Count enabled/disabled
	enabledCount := 0
	disabledCount := 0

	for _, rule := range firewallRules {
		actionCounts[rule.Action]++
		protocolCounts[rule.Protocol]++
		if rule.Enabled {
			enabledCount++
		} else {
			disabledCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":       len(firewallRules),
		"enabled":     enabledCount,
		"disabled":    disabledCount,
		"by_action":   actionCounts,
		"by_protocol": protocolCounts,
	})
}

// getSystemAnalytics returns overall system analytics
func getSystemAnalytics(c *gin.Context) {
	dataMux.RLock()
	alertCount := len(alerts)
	threatCount := len(threats)
	ruleCount := len(firewallRules)
	dataMux.RUnlock()

	usersMux.RLock()
	userCount := len(users)
	usersMux.RUnlock()

	activityMux.RLock()
	activityCount := len(activities)
	activityMux.RUnlock()

	notificationsMux.RLock()
	notifCount := 0
	for _, notifs := range notifications {
		notifCount += len(notifs)
	}
	notificationsMux.RUnlock()

	// Calculate trends (mock data for now)
	trends := gin.H{
		"alerts_trend":   "+15%",
		"threats_trend":  "+8%",
		"users_trend":    "+3%",
		"activity_trend": "+25%",
	}

	c.JSON(http.StatusOK, gin.H{
		"totals": gin.H{
			"alerts":        alertCount,
			"threats":       threatCount,
			"firewall_rules": ruleCount,
			"users":         userCount,
			"activities":    activityCount,
			"notifications": notifCount,
		},
		"trends": trends,
		"timestamp": time.Now(),
	})
}

// getTimeSeriesData returns time series data for charts
func getTimeSeriesData(c *gin.Context) {
	resource := c.Param("resource")
	hours := 24

	dataMux.RLock()
	defer dataMux.RUnlock()

	// Generate time series data
	data := make([]gin.H, hours)
	now := time.Now()

	for i := 0; i < hours; i++ {
		timestamp := now.Add(time.Duration(-hours+i) * time.Hour)
		count := 0

		switch resource {
		case "alerts":
			for _, alert := range alerts {
				if alert.Timestamp.Hour() == timestamp.Hour() {
					count++
				}
			}
		case "threats":
			for _, threat := range threats {
				if threat.Timestamp.Hour() == timestamp.Hour() {
					count++
				}
			}
		}

		data[i] = gin.H{
			"timestamp": timestamp.Format("2006-01-02 15:00"),
			"count":     count,
			"hour":      timestamp.Hour(),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"resource": resource,
		"period":   "24h",
		"data":     data,
	})
}
