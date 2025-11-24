package main

import (
	"sync"
	"time"
)

// Activity represents a system activity log
type Activity struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	UserEmail   string                 `json:"user_email"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	ResourceID  string                 `json:"resource_id"`
	Details     map[string]interface{} `json:"details"`
	IPAddress   string                 `json:"ip_address"`
	Timestamp   time.Time              `json:"timestamp"`
	Status      string                 `json:"status"`
}

var (
	activities    = make(map[string]*Activity)
	activityMux   sync.RWMutex
	activityCount int
)

// logActivity logs a user activity
func logActivity(userID, userEmail, action, resource, resourceID, ipAddress, status string, details map[string]interface{}) {
	activityMux.Lock()
	defer activityMux.Unlock()

	activityCount++
	id := time.Now().Format("20060102150405") + "-" + string(rune(activityCount))

	activity := &Activity{
		ID:         id,
		UserID:     userID,
		UserEmail:  userEmail,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		Timestamp:  time.Now(),
		Status:     status,
	}

	activities[id] = activity

	// Keep only last 1000 activities
	if len(activities) > 1000 {
		// Find oldest and delete
		var oldestID string
		var oldestTime time.Time
		for id, act := range activities {
			if oldestTime.IsZero() || act.Timestamp.Before(oldestTime) {
				oldestTime = act.Timestamp
				oldestID = id
			}
		}
		delete(activities, oldestID)
	}
}

// getActivities returns recent activities
func getActivities(limit int, userID string) []*Activity {
	activityMux.RLock()
	defer activityMux.RUnlock()

	var result []*Activity
	for _, activity := range activities {
		if userID == "" || activity.UserID == userID {
			result = append(result, activity)
		}
	}

	// Sort by timestamp (newest first)
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Timestamp.Before(result[j].Timestamp) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// Limit results
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}
