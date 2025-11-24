package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Webhook represents a webhook configuration
type Webhook struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	URL       string    `json:"url"`
	Events    []string  `json:"events"`
	Secret    string    `json:"secret"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	webhooks    = make(map[string]*Webhook)
	webhooksMux sync.RWMutex
)

// createWebhook creates a new webhook
func createWebhook(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		URL    string   `json:"url" binding:"required,url"`
		Events []string `json:"events" binding:"required"`
		Secret string   `json:"secret"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := fmt.Sprintf("wh_%d", time.Now().Unix())

	webhook := &Webhook{
		ID:        id,
		UserID:    userID.(string),
		URL:       req.URL,
		Events:    req.Events,
		Secret:    req.Secret,
		Enabled:   true,
		CreatedAt: time.Now(),
	}

	webhooksMux.Lock()
	webhooks[id] = webhook
	webhooksMux.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Webhook created successfully",
		"webhook": webhook,
	})
}

// listWebhooks lists user's webhooks
func listWebhooks(c *gin.Context) {
	userID, _ := c.Get("user_id")

	webhooksMux.RLock()
	defer webhooksMux.RUnlock()

	var userWebhooks []*Webhook
	for _, webhook := range webhooks {
		if webhook.UserID == userID.(string) {
			userWebhooks = append(userWebhooks, webhook)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"webhooks": userWebhooks,
		"total":    len(userWebhooks),
	})
}

// deleteWebhook deletes a webhook
func deleteWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	userID, _ := c.Get("user_id")

	webhooksMux.Lock()
	defer webhooksMux.Unlock()

	if webhook, exists := webhooks[webhookID]; exists {
		if webhook.UserID == userID.(string) {
			delete(webhooks, webhookID)
			c.JSON(http.StatusOK, gin.H{
				"message": "Webhook deleted successfully",
				"id":      webhookID,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
}

// updateWebhook updates a webhook
func updateWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		URL     string   `json:"url"`
		Events  []string `json:"events"`
		Enabled *bool    `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	webhooksMux.Lock()
	defer webhooksMux.Unlock()

	if webhook, exists := webhooks[webhookID]; exists {
		if webhook.UserID == userID.(string) {
			if req.URL != "" {
				webhook.URL = req.URL
			}
			if len(req.Events) > 0 {
				webhook.Events = req.Events
			}
			if req.Enabled != nil {
				webhook.Enabled = *req.Enabled
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Webhook updated successfully",
				"webhook": webhook,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
}

// triggerWebhook sends webhook notification
func triggerWebhook(event string, data interface{}) {
	webhooksMux.RLock()
	defer webhooksMux.RUnlock()

	for _, webhook := range webhooks {
		if !webhook.Enabled {
			continue
		}

		// Check if webhook subscribes to this event
		subscribed := false
		for _, e := range webhook.Events {
			if e == event || e == "*" {
				subscribed = true
				break
			}
		}

		if !subscribed {
			continue
		}

		// Send webhook in goroutine
		go sendWebhook(webhook, event, data)
	}
}

// sendWebhook sends HTTP POST to webhook URL
func sendWebhook(webhook *Webhook, event string, data interface{}) {
	payload := map[string]interface{}{
		"event":     event,
		"data":      data,
		"timestamp": time.Now().Format(time.RFC3339),
		"webhook_id": webhook.ID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "NetGuard-Webhook/1.0")
	if webhook.Secret != "" {
		req.Header.Set("X-Webhook-Secret", webhook.Secret)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// testWebhook sends a test webhook
func testWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	userID, _ := c.Get("user_id")

	webhooksMux.RLock()
	webhook, exists := webhooks[webhookID]
	webhooksMux.RUnlock()

	if !exists || webhook.UserID != userID.(string) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
		return
	}

	// Send test event
	go sendWebhook(webhook, "test", gin.H{
		"message": "This is a test webhook",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Test webhook sent",
		"url":     webhook.URL,
	})
}
