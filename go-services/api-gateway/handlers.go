package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Alert represents a security alert
type Alert struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Source      string    `json:"source"`
}

// Threat represents a security threat
type Threat struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Severity   string    `json:"severity"`
	Status     string    `json:"status"`
	SourceIP   string    `json:"source_ip"`
	TargetIP   string    `json:"target_ip"`
	Port       int       `json:"port"`
	Timestamp  time.Time `json:"timestamp"`
	Detections int       `json:"detections"`
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Action    string    `json:"action"`
	Protocol  string    `json:"protocol"`
	SourceIP  string    `json:"source_ip"`
	DestIP    string    `json:"dest_ip"`
	Port      int       `json:"port"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name      string `json:"name"`
	IP        string `json:"ip"`
	MAC       string `json:"mac"`
	Status    string `json:"status"`
	BytesSent int64  `json:"bytes_sent"`
	BytesRecv int64  `json:"bytes_recv"`
}

// In-memory storage
var (
	alerts        = make(map[string]*Alert)
	threats       = make(map[string]*Threat)
	firewallRules = make(map[string]*FirewallRule)
	dataMux       sync.RWMutex
)

func init() {
	// Initialize with sample data
	initSampleData()
}

func initSampleData() {
	// Sample alerts
	alerts["ALT-001"] = &Alert{
		ID:          "ALT-001",
		Title:       "Suspicious Login Attempt",
		Description: "Multiple failed login attempts detected from IP 192.168.1.100",
		Severity:    "critical",
		Status:      "active",
		Timestamp:   time.Now(),
		Source:      "Authentication System",
	}
	alerts["ALT-002"] = &Alert{
		ID:          "ALT-002",
		Title:       "Unusual Network Traffic",
		Description: "High volume of outbound traffic detected",
		Severity:    "high",
		Status:      "investigating",
		Timestamp:   time.Now().Add(-15 * time.Minute),
		Source:      "Network Monitor",
	}
	alerts["ALT-003"] = &Alert{
		ID:          "ALT-003",
		Title:       "Port Scan Detected",
		Description: "Port scanning activity from external IP",
		Severity:    "high",
		Status:      "active",
		Timestamp:   time.Now().Add(-30 * time.Minute),
		Source:      "IDS",
	}

	// Sample threats
	threats["THR-001"] = &Threat{
		ID:         "THR-001",
		Name:       "Malware.Generic.Trojan",
		Type:       "Malware",
		Severity:   "critical",
		Status:     "blocked",
		SourceIP:   "192.168.1.100",
		TargetIP:   "10.0.0.1",
		Port:       443,
		Timestamp:  time.Now(),
		Detections: 145,
	}
	threats["THR-002"] = &Threat{
		ID:         "THR-002",
		Name:       "Brute.Force.SSH",
		Type:       "Brute Force",
		Severity:   "high",
		Status:     "monitoring",
		SourceIP:   "203.0.113.45",
		TargetIP:   "10.0.0.5",
		Port:       22,
		Timestamp:  time.Now().Add(-1 * time.Hour),
		Detections: 67,
	}

	// Sample firewall rules
	firewallRules["FW-001"] = &FirewallRule{
		ID:        "FW-001",
		Name:      "Block Malicious IP",
		Action:    "deny",
		Protocol:  "tcp",
		SourceIP:  "192.168.1.100",
		DestIP:    "any",
		Port:      0,
		Enabled:   true,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}
	firewallRules["FW-002"] = &FirewallRule{
		ID:        "FW-002",
		Name:      "Allow HTTPS",
		Action:    "allow",
		Protocol:  "tcp",
		SourceIP:  "any",
		DestIP:    "any",
		Port:      443,
		Enabled:   true,
		CreatedAt: time.Now().Add(-48 * time.Hour),
	}
	firewallRules["FW-003"] = &FirewallRule{
		ID:        "FW-003",
		Name:      "Allow SSH",
		Action:    "allow",
		Protocol:  "tcp",
		SourceIP:  "10.0.0.0/24",
		DestIP:    "any",
		Port:      22,
		Enabled:   true,
		CreatedAt: time.Now().Add(-72 * time.Hour),
	}
}

// Alert handlers
func listAlerts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	severity := c.Query("severity")

	dataMux.RLock()
	defer dataMux.RUnlock()

	var alertList []*Alert
	for _, alert := range alerts {
		if severity != "" && alert.Severity != severity {
			continue
		}
		alertList = append(alertList, alert)
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alertList,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(alertList),
		},
		"filters": gin.H{
			"severity": severity,
		},
	})
}

func getAlert(c *gin.Context) {
	id := c.Param("id")

	dataMux.RLock()
	alert, exists := alerts[id]
	dataMux.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	c.JSON(http.StatusOK, alert)
}

func createAlert(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Severity    string `json:"severity" binding:"required"`
		Source      string `json:"source"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := fmt.Sprintf("ALT-%03d", len(alerts)+1)
	alert := &Alert{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Severity:    req.Severity,
		Status:      "active",
		Timestamp:   time.Now(),
		Source:      req.Source,
	}

	dataMux.Lock()
	alerts[id] = alert
	dataMux.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Alert created successfully",
		"alert":   alert,
	})
}

func updateAlert(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dataMux.Lock()
	alert, exists := alerts[id]
	if exists {
		alert.Status = req.Status
	}
	dataMux.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Alert updated successfully",
		"alert":   alert,
	})
}

func deleteAlert(c *gin.Context) {
	id := c.Param("id")

	dataMux.Lock()
	_, exists := alerts[id]
	if exists {
		delete(alerts, id)
	}
	dataMux.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Alert deleted successfully",
		"id":      id,
	})
}

// All other handlers are now in handlers_impl.go
