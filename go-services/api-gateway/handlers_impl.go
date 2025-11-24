package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Threat handlers
func listThreats(c *gin.Context) {
	dataMux.RLock()
	defer dataMux.RUnlock()

	var threatList []*Threat
	for _, threat := range threats {
		threatList = append(threatList, threat)
	}

	c.JSON(http.StatusOK, gin.H{
		"threats": threatList,
		"total":   len(threatList),
	})
}

func getThreat(c *gin.Context) {
	id := c.Param("id")

	dataMux.RLock()
	threat, exists := threats[id]
	dataMux.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Threat not found"})
		return
	}

	c.JSON(http.StatusOK, threat)
}

func analyzeThreat(c *gin.Context) {
	var req struct {
		Data string `json:"data" binding:"required"`
		Type string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create new threat from analysis
	id := fmt.Sprintf("THR-%03d", len(threats)+1)
	threat := &Threat{
		ID:         id,
		Name:       fmt.Sprintf("Analyzed.%s", req.Type),
		Type:       req.Type,
		Severity:   "medium",
		Status:     "analyzing",
		SourceIP:   "unknown",
		TargetIP:   "unknown",
		Port:       0,
		Timestamp:  time.Now(),
		Detections: 1,
	}

	dataMux.Lock()
	threats[id] = threat
	dataMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"analysis_id": id,
		"status":      "completed",
		"threat":      threat,
		"message":     "Threat analysis completed",
	})
}

// Firewall handlers
func listFirewallRules(c *gin.Context) {
	dataMux.RLock()
	defer dataMux.RUnlock()

	var ruleList []*FirewallRule
	for _, rule := range firewallRules {
		ruleList = append(ruleList, rule)
	}

	c.JSON(http.StatusOK, gin.H{
		"rules": ruleList,
		"total": len(ruleList),
	})
}

func addFirewallRule(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Action   string `json:"action" binding:"required"`
		Protocol string `json:"protocol" binding:"required"`
		SourceIP string `json:"source_ip"`
		DestIP   string `json:"dest_ip"`
		Port     int    `json:"port"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := fmt.Sprintf("FW-%03d", len(firewallRules)+1)
	rule := &FirewallRule{
		ID:        id,
		Name:      req.Name,
		Action:    req.Action,
		Protocol:  req.Protocol,
		SourceIP:  req.SourceIP,
		DestIP:    req.DestIP,
		Port:      req.Port,
		Enabled:   true,
		CreatedAt: time.Now(),
	}

	dataMux.Lock()
	firewallRules[id] = rule
	dataMux.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Firewall rule added successfully",
		"rule":    rule,
	})
}

func deleteFirewallRule(c *gin.Context) {
	id := c.Param("id")

	dataMux.Lock()
	_, exists := firewallRules[id]
	if exists {
		delete(firewallRules, id)
	}
	dataMux.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Firewall rule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Firewall rule deleted successfully",
		"id":      id,
	})
}

// Network handlers
func listInterfaces(c *gin.Context) {
	interfaces := []NetworkInterface{
		{
			Name:      "eth0",
			IP:        "192.168.1.10",
			MAC:       "00:11:22:33:44:55",
			Status:    "up",
			BytesSent: 1024000000,
			BytesRecv: 2048000000,
		},
		{
			Name:      "eth1",
			IP:        "10.0.0.5",
			MAC:       "AA:BB:CC:DD:EE:FF",
			Status:    "up",
			BytesSent: 512000000,
			BytesRecv: 1024000000,
		},
		{
			Name:      "lo",
			IP:        "127.0.0.1",
			MAC:       "00:00:00:00:00:00",
			Status:    "up",
			BytesSent: 1024000,
			BytesRecv: 1024000,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"interfaces": interfaces,
		"total":      len(interfaces),
	})
}

func getNetworkStats(c *gin.Context) {
	stats := gin.H{
		"packets_captured":  15234 + (time.Now().Unix() % 1000),
		"bytes_processed":   1024000000 + (time.Now().Unix() % 1000000),
		"alerts_generated":  42 + (time.Now().Unix() % 10),
		"threats_detected":  7 + (time.Now().Unix() % 5),
		"uptime_seconds":    time.Now().Unix() % 86400,
		"active_connections": 42 + (time.Now().Second() % 20),
		"bandwidth_usage": gin.H{
			"inbound":  "125.5 MB/s",
			"outbound": "87.3 MB/s",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

func startMonitoring(c *gin.Context) {
	var req struct {
		Interface string `json:"interface" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Monitoring started successfully",
		"interface": req.Interface,
		"status":    "active",
		"started_at": time.Now(),
	})
}

func stopMonitoring(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":    "Monitoring stopped successfully",
		"status":     "stopped",
		"stopped_at": time.Now(),
	})
}

// User management handlers
func listUsers(c *gin.Context) {
	usersMux.RLock()
	defer usersMux.RUnlock()

	var userList []gin.H
	for _, user := range users {
		userList = append(userList, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"name":       user.Name,
			"company":    user.Company,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
		"total": len(userList),
	})
}

func getUser(c *gin.Context) {
	id := c.Param("id")

	usersMux.RLock()
	defer usersMux.RUnlock()

	var foundUser *User
	for _, user := range users {
		if user.ID == id {
			foundUser = user
			break
		}
	}

	if foundUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         foundUser.ID,
		"email":      foundUser.Email,
		"name":       foundUser.Name,
		"company":    foundUser.Company,
		"role":       foundUser.Role,
		"created_at": foundUser.CreatedAt,
	})
}

func updateUser(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name    string `json:"name"`
		Company string `json:"company"`
		Role    string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersMux.Lock()
	var foundUser *User
	for _, user := range users {
		if user.ID == id {
			foundUser = user
			if req.Name != "" {
				user.Name = req.Name
			}
			if req.Company != "" {
				user.Company = req.Company
			}
			if req.Role != "" {
				user.Role = req.Role
			}
			break
		}
	}
	usersMux.Unlock()

	if foundUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "User updated successfully",
		"user": gin.H{
			"id":      foundUser.ID,
			"email":   foundUser.Email,
			"name":    foundUser.Name,
			"company": foundUser.Company,
			"role":    foundUser.Role,
		},
	})
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	usersMux.Lock()
	var found bool
	for email, user := range users {
		if user.ID == id {
			delete(users, email)
			found = true
			break
		}
	}
	usersMux.Unlock()

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// Dashboard handlers
func getDashboardStats(c *gin.Context) {
	dataMux.RLock()
	alertCount := len(alerts)
	threatCount := len(threats)
	firewallRuleCount := len(firewallRules)
	dataMux.RUnlock()

	usersMux.RLock()
	userCount := len(users)
	usersMux.RUnlock()

	stats := gin.H{
		"total_alerts":        alertCount,
		"active_alerts":       alertCount,
		"total_threats":       threatCount,
		"blocked_threats":     threatCount / 2,
		"firewall_rules":      firewallRuleCount,
		"active_users":        userCount,
		"network_interfaces":  3,
		"packets_per_second":  1234.5 + float64(time.Now().Second()%100),
		"bytes_per_second":    5242880 + (time.Now().Unix() % 1000000),
		"active_connections":  42 + (time.Now().Second() % 20),
		"threats_detected":    7 + (time.Now().Second() % 10),
		"uptime_hours":        24,
		"cpu_usage":           45.5,
		"memory_usage":        62.3,
		"disk_usage":          78.9,
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}

func getRecentActivity(c *gin.Context) {
	dataMux.RLock()
	defer dataMux.RUnlock()

	var activities []gin.H

	// Add recent alerts
	for _, alert := range alerts {
		activities = append(activities, gin.H{
			"id":        alert.ID,
			"type":      "alert",
			"title":     alert.Title,
			"severity":  alert.Severity,
			"timestamp": alert.Timestamp,
		})
	}

	// Add recent threats
	for _, threat := range threats {
		activities = append(activities, gin.H{
			"id":        threat.ID,
			"type":      "threat",
			"title":     threat.Name,
			"severity":  threat.Severity,
			"timestamp": threat.Timestamp,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}
