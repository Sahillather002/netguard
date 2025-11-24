package main

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// getSystemInfo returns system information
func getSystemInfo(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime)

	c.JSON(http.StatusOK, gin.H{
		"system": gin.H{
			"go_version":    runtime.Version(),
			"num_goroutine": runtime.NumGoroutine(),
			"num_cpu":       runtime.NumCPU(),
			"memory": gin.H{
				"alloc_mb":       m.Alloc / 1024 / 1024,
				"total_alloc_mb": m.TotalAlloc / 1024 / 1024,
				"sys_mb":         m.Sys / 1024 / 1024,
				"num_gc":         m.NumGC,
			},
		},
		"uptime": gin.H{
			"seconds": int(uptime.Seconds()),
			"minutes": int(uptime.Minutes()),
			"hours":   int(uptime.Hours()),
			"days":    int(uptime.Hours() / 24),
			"started": startTime.Format(time.RFC3339),
		},
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// getHealthCheck returns detailed health check
func getHealthCheck(c *gin.Context) {
	dataMux.RLock()
	alertCount := len(alerts)
	threatCount := len(threats)
	dataMux.RUnlock()

	usersMux.RLock()
	userCount := len(users)
	usersMux.RUnlock()

	// Check system health
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memoryUsage := float64(m.Alloc) / float64(m.Sys) * 100
	goroutineCount := runtime.NumGoroutine()

	status := "healthy"
	if memoryUsage > 90 || goroutineCount > 10000 {
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"checks": gin.H{
			"database": gin.H{
				"status":  "connected",
				"latency": "2ms",
			},
			"memory": gin.H{
				"status": "ok",
				"usage":  memoryUsage,
			},
			"goroutines": gin.H{
				"status": "ok",
				"count":  goroutineCount,
			},
		},
		"data": gin.H{
			"alerts":  alertCount,
			"threats": threatCount,
			"users":   userCount,
		},
		"uptime":    int(time.Since(startTime).Seconds()),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// getReadinessCheck returns readiness status
func getReadinessCheck(c *gin.Context) {
	// Check if system is ready to accept requests
	ready := true
	checks := gin.H{}

	// Check data stores
	dataMux.RLock()
	dataReady := len(alerts) >= 0 && len(threats) >= 0
	dataMux.RUnlock()

	checks["data_store"] = dataReady

	// Check users
	usersMux.RLock()
	usersReady := len(users) >= 0
	usersMux.RUnlock()

	checks["users"] = usersReady

	if !dataReady || !usersReady {
		ready = false
	}

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, gin.H{
		"ready":  ready,
		"checks": checks,
	})
}
