package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID            string                 `json:"id"`
	Timestamp     time.Time              `json:"timestamp"`
	UserID        string                 `json:"user_id"`
	UserEmail     string                 `json:"user_email"`
	Action        string                 `json:"action"`
	Resource      string                 `json:"resource"`
	ResourceID    string                 `json:"resource_id"`
	Method        string                 `json:"method"`
	Path          string                 `json:"path"`
	IPAddress     string                 `json:"ip_address"`
	UserAgent     string                 `json:"user_agent"`
	StatusCode    int                    `json:"status_code"`
	RequestBody   string                 `json:"request_body,omitempty"`
	ResponseBody  string                 `json:"response_body,omitempty"`
	Duration      int64                  `json:"duration_ms"`
	Changes       map[string]interface{} `json:"changes,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

var (
	auditLogs    = make([]*AuditLog, 0)
	auditLogsMux sync.RWMutex
	auditLogID   int
)

// auditMiddleware logs all requests for audit trail
func auditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Capture request body
		var requestBody string
		if c.Request.Body != nil && c.Request.Method != "GET" {
			bodyBytes, _ := c.GetRawData()
			requestBody = string(bodyBytes)
			// Restore body for next handlers
			c.Request.Body = http.MaxBytesReader(c.Writer, http.NoBody, 0)
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime).Milliseconds()

		// Get user info
		userID, _ := c.Get("user_id")
		userEmail := ""
		if user, exists := c.Get("user"); exists {
			if u, ok := user.(*User); ok {
				userEmail = u.Email
			}
		}

		// Create audit log
		auditLogsMux.Lock()
		auditLogID++
		logID := fmt.Sprintf("audit_%d", auditLogID)

		log := &AuditLog{
			ID:           logID,
			Timestamp:    startTime,
			UserID:       fmt.Sprintf("%v", userID),
			UserEmail:    userEmail,
			Action:       c.Request.Method,
			Resource:     c.FullPath(),
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			StatusCode:   c.Writer.Status(),
			RequestBody:  requestBody,
			Duration:     duration,
		}

		auditLogs = append(auditLogs, log)

		// Keep only last 10000 logs
		if len(auditLogs) > 10000 {
			auditLogs = auditLogs[len(auditLogs)-10000:]
		}

		auditLogsMux.Unlock()
	}
}

// getAuditLogs returns audit logs with filtering
func getAuditLogs(c *gin.Context) {
	userID := c.Query("user_id")
	action := c.Query("action")
	resource := c.Query("resource")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	auditLogsMux.RLock()
	defer auditLogsMux.RUnlock()

	var filtered []*AuditLog
	for _, log := range auditLogs {
		// Apply filters
		if userID != "" && log.UserID != userID {
			continue
		}
		if action != "" && log.Action != action {
			continue
		}
		if resource != "" && log.Resource != resource {
			continue
		}
		if startDate != "" {
			start, _ := time.Parse(time.RFC3339, startDate)
			if log.Timestamp.Before(start) {
				continue
			}
		}
		if endDate != "" {
			end, _ := time.Parse(time.RFC3339, endDate)
			if log.Timestamp.After(end) {
				continue
			}
		}

		filtered = append(filtered, log)
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  filtered,
		"total": len(filtered),
	})
}

// getAuditLog returns a single audit log
func getAuditLog(c *gin.Context) {
	logID := c.Param("id")

	auditLogsMux.RLock()
	defer auditLogsMux.RUnlock()

	for _, log := range auditLogs {
		if log.ID == logID {
			c.JSON(http.StatusOK, log)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
}

// exportAuditLogs exports audit logs
func exportAuditLogs(c *gin.Context) {
	format := c.DefaultQuery("format", "json")

	auditLogsMux.RLock()
	logs := make([]*AuditLog, len(auditLogs))
	copy(logs, auditLogs)
	auditLogsMux.RUnlock()

	switch format {
	case "json":
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=audit_logs_%s.json", time.Now().Format("20060102")))
		
		data, _ := json.MarshalIndent(gin.H{
			"logs":        logs,
			"total":       len(logs),
			"exported_at": time.Now(),
		}, "", "  ")
		
		c.Data(http.StatusOK, "application/json", data)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
	}
}

// getAuditStats returns audit statistics
func getAuditStats(c *gin.Context) {
	auditLogsMux.RLock()
	defer auditLogsMux.RUnlock()

	// Count by action
	actionCounts := make(map[string]int)
	// Count by resource
	resourceCounts := make(map[string]int)
	// Count by user
	userCounts := make(map[string]int)
	// Count by status code
	statusCounts := make(map[int]int)
	// Average duration
	var totalDuration int64

	for _, log := range auditLogs {
		actionCounts[log.Action]++
		resourceCounts[log.Resource]++
		userCounts[log.UserID]++
		statusCounts[log.StatusCode]++
		totalDuration += log.Duration
	}

	avgDuration := int64(0)
	if len(auditLogs) > 0 {
		avgDuration = totalDuration / int64(len(auditLogs))
	}

	c.JSON(http.StatusOK, gin.H{
		"total_logs":       len(auditLogs),
		"by_action":        actionCounts,
		"by_resource":      resourceCounts,
		"by_user":          userCounts,
		"by_status":        statusCounts,
		"avg_duration_ms":  avgDuration,
	})
}
