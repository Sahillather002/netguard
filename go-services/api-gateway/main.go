package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Initialize sample notifications
	initNotifications()

	// Start rate limit cleanup
	startRateLimitCleanup()

	// Start cache cleanup
	startCacheCleanup()

	// Initialize Gin router
	router := gin.Default()

	// Add middlewares
	router.Use(rateLimitMiddleware())
	router.Use(performanceMiddleware())

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:4200"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "X-API-Key", "Accept"}
	config.ExposeHeaders = []string{"Content-Length", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Health check endpoints
	router.GET("/health", getHealthCheck)
	router.GET("/ready", getReadinessCheck)
	router.GET("/system/info", getSystemInfo)

	// Metrics endpoint for Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", login)
			auth.POST("/register", register)
			auth.POST("/refresh", refreshToken)
			auth.POST("/logout", logout)
		}

		// Protected routes (require authentication or API key)
		protected := v1.Group("/")
		protected.Use(authOrAPIKeyMiddleware())
		protected.Use(loggingMiddleware())
		protected.Use(auditMiddleware())
		{
			// Rate limit status
			protected.GET("/ratelimit", getRateLimitStatus)
			// User profile
			protected.GET("/me", func(c *gin.Context) {
				user, _ := c.Get("user")
				c.JSON(http.StatusOK, gin.H{"user": user})
			})

			// Search
			protected.GET("/search", searchAll)

			// Activity logs
			protected.GET("/activities", listActivities)
			protected.GET("/activities/stats", getActivityStats)

			// Export endpoints
			protected.GET("/export/alerts", exportAlerts)
			protected.GET("/export/threats", exportThreats)
			protected.GET("/export/firewall-rules", exportFirewallRules)

			// Notifications
			protected.GET("/notifications", listNotifications)
			protected.PUT("/notifications/:id/read", markNotificationRead)
			protected.POST("/notifications/read-all", markAllNotificationsRead)
			protected.DELETE("/notifications/:id", deleteNotification)

			// Analytics
			protected.GET("/analytics/alerts", getAlertAnalytics)
			protected.GET("/analytics/threats", getThreatAnalytics)
			protected.GET("/analytics/firewall", getFirewallAnalytics)
			protected.GET("/analytics/system", getSystemAnalytics)
			protected.GET("/analytics/timeseries/:resource", getTimeSeriesData)

			// Batch operations
			protected.POST("/batch/alerts/delete", batchDeleteAlerts)
			protected.POST("/batch/alerts/update", batchUpdateAlerts)
			protected.POST("/batch/threats/delete", batchDeleteThreats)
			protected.POST("/batch/firewall-rules/delete", batchDeleteFirewallRules)
			protected.POST("/batch/firewall-rules/enable", batchEnableFirewallRules)

			// API Keys
			protected.GET("/api-keys", listAPIKeys)
			protected.POST("/api-keys", createAPIKey)
			protected.PUT("/api-keys/:id", updateAPIKey)
			protected.DELETE("/api-keys/:id", revokeAPIKey)

			// Webhooks
			protected.GET("/webhooks", listWebhooks)
			protected.POST("/webhooks", createWebhook)
			protected.PUT("/webhooks/:id", updateWebhook)
			protected.DELETE("/webhooks/:id", deleteWebhook)
			protected.POST("/webhooks/:id/test", testWebhook)

			// Audit Logs
			protected.GET("/audit-logs", getAuditLogs)
			protected.GET("/audit-logs/:id", getAuditLog)
			protected.GET("/audit-logs/export", exportAuditLogs)
			protected.GET("/audit-logs/stats", getAuditStats)

			// Compliance
			protected.GET("/compliance/report", generateComplianceReport)
			protected.GET("/compliance/status", getComplianceStatus)
			protected.GET("/compliance/checklist", getComplianceChecklist)

			// Reports
			protected.GET("/reports/security", generateSecurityReport)
			protected.GET("/reports/threats", generateThreatReport)
			protected.GET("/reports/network", generateNetworkReport)
			protected.POST("/reports/schedule", scheduleReport)

			// Performance
			protected.GET("/performance/metrics", getPerformanceMetrics)
			protected.GET("/performance/slowest", getSlowestEndpoints)
			protected.GET("/performance/most-used", getMostUsedEndpoints)
			protected.POST("/performance/reset", resetPerformanceMetrics)

			// Backup & Restore
			protected.POST("/backup/create", createBackup)
			protected.GET("/backup/download", downloadBackup)
			protected.POST("/backup/restore", restoreBackup)
			protected.GET("/backup/info", getBackupInfo)

			// Cache
			protected.GET("/cache/stats", func(c *gin.Context) {
				c.JSON(http.StatusOK, cache.GetStats())
			})
			protected.POST("/cache/clear", func(c *gin.Context) {
				cache.Clear()
				c.JSON(http.StatusOK, gin.H{"message": "Cache cleared"})
			})

			// Alerts endpoints
			alerts := protected.Group("/alerts")
			{
				alerts.GET("", listAlerts)
				alerts.GET("/:id", getAlert)
				alerts.POST("", createAlert)
				alerts.PUT("/:id", updateAlert)
				alerts.DELETE("/:id", deleteAlert)
			}

			// Network monitoring endpoints
			network := protected.Group("/network")
			{
				network.GET("/interfaces", listInterfaces)
				network.GET("/stats", getNetworkStats)
				network.POST("/monitor/start", startMonitoring)
				network.POST("/monitor/stop", stopMonitoring)
			}

			// Firewall rules endpoints
			firewall := protected.Group("/firewall")
			{
				firewall.GET("/rules", listFirewallRules)
				firewall.POST("/rules", addFirewallRule)
				firewall.DELETE("/rules/:id", deleteFirewallRule)
			}

			// Threat detection endpoints
			threats := protected.Group("/threats")
			{
				threats.GET("", listThreats)
				threats.GET("/:id", getThreat)
				threats.POST("/analyze", analyzeThreat)
			}

			// User management endpoints
			users := protected.Group("/users")
			{
				users.GET("", listUsers)
				users.GET("/:id", getUser)
				users.PUT("/:id", updateUser)
				users.DELETE("/:id", deleteUser)
			}

			// Dashboard endpoints
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", getDashboardStats)
				dashboard.GET("/recent-activity", getRecentActivity)
			}
		}
	}

	// GraphQL endpoint
	router.POST("/graphql", graphqlHandler)
	router.GET("/graphql", playgroundHandler)

	// WebSocket endpoints for real-time updates
	router.GET("/ws", websocketHandler)
	router.GET("/ws/stats", statsWebSocketHandler)
	router.GET("/ws/threats", threatsWebSocketHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("🚀 API Gateway starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server exited")
}
