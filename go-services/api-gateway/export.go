package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// exportAlerts exports alerts to CSV or JSON
func exportAlerts(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	dataMux.RLock()
	var alertList []*Alert
	for _, alert := range alerts {
		alertList = append(alertList, alert)
	}
	dataMux.RUnlock()

	switch format {
	case "csv":
		exportAlertsCSV(c, alertList)
	case "json":
		exportAlertsJSON(c, alertList)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Use 'csv' or 'json'"})
	}
}

func exportAlertsCSV(c *gin.Context, alerts []*Alert) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=alerts_%s.csv", time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Title", "Description", "Severity", "Status", "Source", "Timestamp"})

	// Write data
	for _, alert := range alerts {
		writer.Write([]string{
			alert.ID,
			alert.Title,
			alert.Description,
			alert.Severity,
			alert.Status,
			alert.Source,
			alert.Timestamp.Format(time.RFC3339),
		})
	}
}

func exportAlertsJSON(c *gin.Context, alerts []*Alert) {
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=alerts_%s.json", time.Now().Format("20060102")))

	data, _ := json.MarshalIndent(gin.H{
		"alerts":      alerts,
		"total":       len(alerts),
		"exported_at": time.Now(),
	}, "", "  ")

	c.Data(http.StatusOK, "application/json", data)
}

// exportThreats exports threats to CSV or JSON
func exportThreats(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	dataMux.RLock()
	var threatList []*Threat
	for _, threat := range threats {
		threatList = append(threatList, threat)
	}
	dataMux.RUnlock()

	switch format {
	case "csv":
		exportThreatsCSV(c, threatList)
	case "json":
		exportThreatsJSON(c, threatList)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Use 'csv' or 'json'"})
	}
}

func exportThreatsCSV(c *gin.Context, threats []*Threat) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=threats_%s.csv", time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Name", "Type", "Severity", "Status", "Source IP", "Target IP", "Port", "Detections", "Timestamp"})

	// Write data
	for _, threat := range threats {
		writer.Write([]string{
			threat.ID,
			threat.Name,
			threat.Type,
			threat.Severity,
			threat.Status,
			threat.SourceIP,
			threat.TargetIP,
			fmt.Sprintf("%d", threat.Port),
			fmt.Sprintf("%d", threat.Detections),
			threat.Timestamp.Format(time.RFC3339),
		})
	}
}

func exportThreatsJSON(c *gin.Context, threats []*Threat) {
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=threats_%s.json", time.Now().Format("20060102")))

	data, _ := json.MarshalIndent(gin.H{
		"threats":     threats,
		"total":       len(threats),
		"exported_at": time.Now(),
	}, "", "  ")

	c.Data(http.StatusOK, "application/json", data)
}

// exportFirewallRules exports firewall rules
func exportFirewallRules(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	dataMux.RLock()
	var ruleList []*FirewallRule
	for _, rule := range firewallRules {
		ruleList = append(ruleList, rule)
	}
	dataMux.RUnlock()

	switch format {
	case "csv":
		exportFirewallRulesCSV(c, ruleList)
	case "json":
		exportFirewallRulesJSON(c, ruleList)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format. Use 'csv' or 'json'"})
	}
}

func exportFirewallRulesCSV(c *gin.Context, rules []*FirewallRule) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=firewall_rules_%s.csv", time.Now().Format("20060102")))

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Name", "Action", "Protocol", "Source IP", "Dest IP", "Port", "Enabled", "Created At"})

	// Write data
	for _, rule := range rules {
		writer.Write([]string{
			rule.ID,
			rule.Name,
			rule.Action,
			rule.Protocol,
			rule.SourceIP,
			rule.DestIP,
			fmt.Sprintf("%d", rule.Port),
			fmt.Sprintf("%t", rule.Enabled),
			rule.CreatedAt.Format(time.RFC3339),
		})
	}
}

func exportFirewallRulesJSON(c *gin.Context, rules []*FirewallRule) {
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=firewall_rules_%s.json", time.Now().Format("20060102")))

	data, _ := json.MarshalIndent(gin.H{
		"rules":       rules,
		"total":       len(rules),
		"exported_at": time.Now(),
	}, "", "  ")

	c.Data(http.StatusOK, "application/json", data)
}
