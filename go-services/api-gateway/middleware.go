package main

import (
	"github.com/gin-gonic/gin"
)

// Logging middleware
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user info if authenticated
		userID, _ := c.Get("user_id")
		userEmail := ""
		if user, exists := c.Get("user"); exists {
			if u, ok := user.(*User); ok {
				userEmail = u.Email
			}
		}

		// Log after request completes
		c.Next()

		// Log activity for write operations
		if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" {
			action := c.Request.Method
			resource := c.FullPath()
			status := "success"
			if c.Writer.Status() >= 400 {
				status = "failed"
			}

			if uid, ok := userID.(string); ok {
				logActivity(uid, userEmail, action, resource, "", c.ClientIP(), status, nil)
			}
		}
	}
}

// CORS middleware (already using gin-contrib/cors, but custom implementation here)
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Request ID middleware
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or get request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// TODO: Generate UUID
			requestID = "req_123456"
		}

		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
