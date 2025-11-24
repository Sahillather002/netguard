package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Notification represents a user notification
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Read      bool      `json:"read"`
	Link      string    `json:"link"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	notifications    = make(map[string][]*Notification)
	notificationsMux sync.RWMutex
)

// createNotification creates a new notification for a user
func createNotification(userID, notifType, title, message, severity, link string) {
	notificationsMux.Lock()
	defer notificationsMux.Unlock()

	notif := &Notification{
		ID:        time.Now().Format("20060102150405"),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Severity:  severity,
		Read:      false,
		Link:      link,
		Timestamp: time.Now(),
	}

	notifications[userID] = append(notifications[userID], notif)

	// Keep only last 100 notifications per user
	if len(notifications[userID]) > 100 {
		notifications[userID] = notifications[userID][len(notifications[userID])-100:]
	}
}

// listNotifications returns user notifications
func listNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	unreadOnly := c.Query("unread") == "true"

	notificationsMux.RLock()
	defer notificationsMux.RUnlock()

	var userNotifs []*Notification
	if uid, ok := userID.(string); ok {
		for _, notif := range notifications[uid] {
			if !unreadOnly || !notif.Read {
				userNotifs = append(userNotifs, notif)
			}
		}
	}

	// Count unread
	unreadCount := 0
	for _, notif := range userNotifs {
		if !notif.Read {
			unreadCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": userNotifs,
		"total":         len(userNotifs),
		"unread_count":  unreadCount,
	})
}

// markNotificationRead marks a notification as read
func markNotificationRead(c *gin.Context) {
	notifID := c.Param("id")
	userID, _ := c.Get("user_id")

	notificationsMux.Lock()
	defer notificationsMux.Unlock()

	if uid, ok := userID.(string); ok {
		for _, notif := range notifications[uid] {
			if notif.ID == notifID {
				notif.Read = true
				c.JSON(http.StatusOK, gin.H{
					"message": "Notification marked as read",
					"id":      notifID,
				})
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
}

// markAllNotificationsRead marks all notifications as read
func markAllNotificationsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")

	notificationsMux.Lock()
	defer notificationsMux.Unlock()

	if uid, ok := userID.(string); ok {
		count := 0
		for _, notif := range notifications[uid] {
			if !notif.Read {
				notif.Read = true
				count++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "All notifications marked as read",
			"count":   count,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "No notifications to mark"})
}

// deleteNotification deletes a notification
func deleteNotification(c *gin.Context) {
	notifID := c.Param("id")
	userID, _ := c.Get("user_id")

	notificationsMux.Lock()
	defer notificationsMux.Unlock()

	if uid, ok := userID.(string); ok {
		for i, notif := range notifications[uid] {
			if notif.ID == notifID {
				notifications[uid] = append(notifications[uid][:i], notifications[uid][i+1:]...)
				c.JSON(http.StatusOK, gin.H{
					"message": "Notification deleted",
					"id":      notifID,
				})
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
}

// Initialize some sample notifications
func initNotifications() {
	createNotification("user_001", "alert", "New Critical Alert", "Suspicious login attempt detected", "critical", "/dashboard/alerts/ALT-001")
	createNotification("user_001", "threat", "Threat Detected", "Malware detected on system", "high", "/dashboard/threats/THR-001")
	createNotification("user_001", "system", "System Update", "New security patches available", "info", "/dashboard/settings")
}
