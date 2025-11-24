package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// APIKey represents an API key
type APIKey struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Name        string    `json:"name"`
	UserID      string    `json:"user_id"`
	Permissions []string  `json:"permissions"`
	Enabled     bool      `json:"enabled"`
	LastUsed    time.Time `json:"last_used"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

var (
	apiKeys    = make(map[string]*APIKey)
	apiKeysMux sync.RWMutex
)

// generateAPIKey generates a random API key
func generateAPIKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "sk_" + base64.URLEncoding.EncodeToString(b)
}

// createAPIKey creates a new API key
func createAPIKey(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name        string   `json:"name" binding:"required"`
		Permissions []string `json:"permissions"`
		ExpiresIn   int      `json:"expires_in"` // days
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key := generateAPIKey()
	id := fmt.Sprintf("key_%d", time.Now().Unix())

	expiresAt := time.Now().AddDate(0, 0, 365) // Default 1 year
	if req.ExpiresIn > 0 {
		expiresAt = time.Now().AddDate(0, 0, req.ExpiresIn)
	}

	apiKey := &APIKey{
		ID:          id,
		Key:         key,
		Name:        req.Name,
		UserID:      userID.(string),
		Permissions: req.Permissions,
		Enabled:     true,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}

	apiKeysMux.Lock()
	apiKeys[key] = apiKey
	apiKeysMux.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"id":         id,
		"key":        key,
		"name":       req.Name,
		"expires_at": expiresAt,
		"message":    "API key created successfully. Save this key securely, it won't be shown again.",
	})
}

// listAPIKeys lists user's API keys
func listAPIKeys(c *gin.Context) {
	userID, _ := c.Get("user_id")

	apiKeysMux.RLock()
	defer apiKeysMux.RUnlock()

	var userKeys []gin.H
	for _, key := range apiKeys {
		if key.UserID == userID.(string) {
			userKeys = append(userKeys, gin.H{
				"id":         key.ID,
				"name":       key.Name,
				"enabled":    key.Enabled,
				"last_used":  key.LastUsed,
				"created_at": key.CreatedAt,
				"expires_at": key.ExpiresAt,
				"key_prefix": key.Key[:10] + "...",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"api_keys": userKeys,
		"total":    len(userKeys),
	})
}

// revokeAPIKey revokes an API key
func revokeAPIKey(c *gin.Context) {
	keyID := c.Param("id")
	userID, _ := c.Get("user_id")

	apiKeysMux.Lock()
	defer apiKeysMux.Unlock()

	for key, apiKey := range apiKeys {
		if apiKey.ID == keyID && apiKey.UserID == userID.(string) {
			delete(apiKeys, key)
			c.JSON(http.StatusOK, gin.H{
				"message": "API key revoked successfully",
				"id":      keyID,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
}

// updateAPIKey updates an API key
func updateAPIKey(c *gin.Context) {
	keyID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req struct {
		Name    string `json:"name"`
		Enabled *bool  `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKeysMux.Lock()
	defer apiKeysMux.Unlock()

	for _, apiKey := range apiKeys {
		if apiKey.ID == keyID && apiKey.UserID == userID.(string) {
			if req.Name != "" {
				apiKey.Name = req.Name
			}
			if req.Enabled != nil {
				apiKey.Enabled = *req.Enabled
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "API key updated successfully",
				"id":      keyID,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "API key not found"})
}

// validateAPIKey validates an API key
func validateAPIKey(key string) (*APIKey, bool) {
	apiKeysMux.RLock()
	defer apiKeysMux.RUnlock()

	apiKey, exists := apiKeys[key]
	if !exists {
		return nil, false
	}

	// Check if enabled
	if !apiKey.Enabled {
		return nil, false
	}

	// Check if expired
	if time.Now().After(apiKey.ExpiresAt) {
		return nil, false
	}

	// Update last used
	go func() {
		apiKeysMux.Lock()
		apiKey.LastUsed = time.Now()
		apiKeysMux.Unlock()
	}()

	return apiKey, true
}

// apiKeyMiddleware validates API key from header
func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.Next()
			return
		}

		key, valid := validateAPIKey(apiKey)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired API key"})
			c.Abort()
			return
		}

		// Set user info from API key
		c.Set("user_id", key.UserID)
		c.Set("api_key_id", key.ID)
		c.Set("auth_method", "api_key")

		c.Next()
	}
}
