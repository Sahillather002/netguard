package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// listActivities returns activity logs
func listActivities(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	userID := c.Query("user_id")

	activities := getActivities(limit, userID)

	c.JSON(http.StatusOK, gin.H{
		"activities": activities,
		"total":      len(activities),
	})
}

// getActivityStats returns activity statistics
func getActivityStats(c *gin.Context) {
	activityMux.RLock()
	defer activityMux.RUnlock()

	// Count by action
	actionCounts := make(map[string]int)
	// Count by resource
	resourceCounts := make(map[string]int)
	// Count by status
	statusCounts := make(map[string]int)

	for _, activity := range activities {
		actionCounts[activity.Action]++
		resourceCounts[activity.Resource]++
		statusCounts[activity.Status]++
	}

	c.JSON(http.StatusOK, gin.H{
		"total_activities": len(activities),
		"by_action":        actionCounts,
		"by_resource":      resourceCounts,
		"by_status":        statusCounts,
	})
}
