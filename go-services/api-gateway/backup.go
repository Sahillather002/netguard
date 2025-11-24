package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Backup represents a system backup
type Backup struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Data      BackupData `json:"data"`
	Size      int       `json:"size"`
}

// BackupData contains all system data
type BackupData struct {
	Alerts        map[string]*Alert        `json:"alerts"`
	Threats       map[string]*Threat       `json:"threats"`
	FirewallRules map[string]*FirewallRule `json:"firewall_rules"`
	Users         map[string]*User         `json:"users"`
	APIKeys       map[string]*APIKey       `json:"api_keys"`
	Webhooks      map[string]*Webhook      `json:"webhooks"`
}

// createBackup creates a backup of all data
func createBackup(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Collect all data
	dataMux.RLock()
	alertsCopy := make(map[string]*Alert)
	for k, v := range alerts {
		alertsCopy[k] = v
	}
	threatsCopy := make(map[string]*Threat)
	for k, v := range threats {
		threatsCopy[k] = v
	}
	rulesCopy := make(map[string]*FirewallRule)
	for k, v := range firewallRules {
		rulesCopy[k] = v
	}
	dataMux.RUnlock()

	usersMux.RLock()
	usersCopy := make(map[string]*User)
	for k, v := range users {
		usersCopy[k] = v
	}
	usersMux.RUnlock()

	apiKeysMux.RLock()
	keysCopy := make(map[string]*APIKey)
	for k, v := range apiKeys {
		keysCopy[k] = v
	}
	apiKeysMux.RUnlock()

	webhooksMux.RLock()
	webhooksCopy := make(map[string]*Webhook)
	for k, v := range webhooks {
		webhooksCopy[k] = v
	}
	webhooksMux.RUnlock()

	backupData := BackupData{
		Alerts:        alertsCopy,
		Threats:       threatsCopy,
		FirewallRules: rulesCopy,
		Users:         usersCopy,
		APIKeys:       keysCopy,
		Webhooks:      webhooksCopy,
	}

	backup := Backup{
		ID:        "backup_" + time.Now().Format("20060102_150405"),
		Timestamp: time.Now(),
		Data:      backupData,
	}

	// Calculate size
	jsonData, _ := json.Marshal(backup)
	backup.Size = len(jsonData)

	// Log activity
	logActivity(fmt.Sprintf("%v", userID), "", "CREATE_BACKUP", "system", backup.ID, c.ClientIP(), "success", nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup created successfully",
		"backup": gin.H{
			"id":        backup.ID,
			"timestamp": backup.Timestamp,
			"size":      backup.Size,
		},
	})
}

// downloadBackup downloads a backup file
func downloadBackup(c *gin.Context) {
	// Collect all data
	dataMux.RLock()
	alertsCopy := make(map[string]*Alert)
	for k, v := range alerts {
		alertsCopy[k] = v
	}
	threatsCopy := make(map[string]*Threat)
	for k, v := range threats {
		threatsCopy[k] = v
	}
	rulesCopy := make(map[string]*FirewallRule)
	for k, v := range firewallRules {
		rulesCopy[k] = v
	}
	dataMux.RUnlock()

	usersMux.RLock()
	usersCopy := make(map[string]*User)
	for k, v := range users {
		// Don't include password hashes in backup
		userCopy := *v
		userCopy.PasswordHash = ""
		usersCopy[k] = &userCopy
	}
	usersMux.RUnlock()

	backupData := BackupData{
		Alerts:        alertsCopy,
		Threats:       threatsCopy,
		FirewallRules: rulesCopy,
		Users:         usersCopy,
	}

	backup := Backup{
		ID:        "backup_" + time.Now().Format("20060102_150405"),
		Timestamp: time.Now(),
		Data:      backupData,
	}

	jsonData, _ := json.MarshalIndent(backup, "", "  ")
	backup.Size = len(jsonData)

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.json", backup.ID))
	c.Data(http.StatusOK, "application/json", jsonData)
}

// restoreBackup restores data from backup
func restoreBackup(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var backup Backup
	if err := c.ShouldBindJSON(&backup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid backup data"})
		return
	}

	// Restore alerts
	dataMux.Lock()
	alerts = backup.Data.Alerts
	threats = backup.Data.Threats
	firewallRules = backup.Data.FirewallRules
	dataMux.Unlock()

	// Restore users (skip if empty to avoid losing current users)
	if len(backup.Data.Users) > 0 {
		usersMux.Lock()
		for k, v := range backup.Data.Users {
			users[k] = v
		}
		usersMux.Unlock()
	}

	// Log activity
	logActivity(fmt.Sprintf("%v", userID), "", "RESTORE_BACKUP", "system", backup.ID, c.ClientIP(), "success", nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup restored successfully",
		"restored": gin.H{
			"alerts":         len(backup.Data.Alerts),
			"threats":        len(backup.Data.Threats),
			"firewall_rules": len(backup.Data.FirewallRules),
			"users":          len(backup.Data.Users),
		},
	})
}

// getBackupInfo returns backup information
func getBackupInfo(c *gin.Context) {
	dataMux.RLock()
	alertCount := len(alerts)
	threatCount := len(threats)
	ruleCount := len(firewallRules)
	dataMux.RUnlock()

	usersMux.RLock()
	userCount := len(users)
	usersMux.RUnlock()

	apiKeysMux.RLock()
	keyCount := len(apiKeys)
	apiKeysMux.RUnlock()

	webhooksMux.RLock()
	webhookCount := len(webhooks)
	webhooksMux.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"current_data": gin.H{
			"alerts":         alertCount,
			"threats":        threatCount,
			"firewall_rules": ruleCount,
			"users":          userCount,
			"api_keys":       keyCount,
			"webhooks":       webhookCount,
		},
		"backup_available": true,
		"last_backup":      nil,
	})
}
