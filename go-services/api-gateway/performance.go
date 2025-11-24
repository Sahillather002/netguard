package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// PerformanceMetric represents a performance metric
type PerformanceMetric struct {
	Endpoint      string
	Method        string
	Count         int64
	TotalDuration int64
	MinDuration   int64
	MaxDuration   int64
	AvgDuration   float64
	ErrorCount    int64
	LastAccessed  time.Time
}

var (
	perfMetrics    = make(map[string]*PerformanceMetric)
	perfMetricsMux sync.RWMutex
)

// performanceMiddleware tracks performance metrics
func performanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime).Milliseconds()

		// Track metrics
		key := c.Request.Method + ":" + c.FullPath()
		trackPerformance(key, c.Request.Method, c.FullPath(), duration, c.Writer.Status())
	}
}

// trackPerformance records performance metrics
func trackPerformance(key, method, endpoint string, duration int64, statusCode int) {
	perfMetricsMux.Lock()
	defer perfMetricsMux.Unlock()

	metric, exists := perfMetrics[key]
	if !exists {
		metric = &PerformanceMetric{
			Endpoint:    endpoint,
			Method:      method,
			MinDuration: duration,
			MaxDuration: duration,
		}
		perfMetrics[key] = metric
	}

	metric.Count++
	metric.TotalDuration += duration
	metric.LastAccessed = time.Now()

	if duration < metric.MinDuration {
		metric.MinDuration = duration
	}
	if duration > metric.MaxDuration {
		metric.MaxDuration = duration
	}

	metric.AvgDuration = float64(metric.TotalDuration) / float64(metric.Count)

	if statusCode >= 400 {
		metric.ErrorCount++
	}
}

// getPerformanceMetrics returns performance metrics
func getPerformanceMetrics(c *gin.Context) {
	perfMetricsMux.RLock()
	defer perfMetricsMux.RUnlock()

	var metrics []*PerformanceMetric
	for _, metric := range perfMetrics {
		metrics = append(metrics, metric)
	}

	// Calculate totals
	var totalRequests int64
	var totalErrors int64
	var totalDuration int64

	for _, metric := range metrics {
		totalRequests += metric.Count
		totalErrors += metric.ErrorCount
		totalDuration += metric.TotalDuration
	}

	avgResponseTime := float64(0)
	if totalRequests > 0 {
		avgResponseTime = float64(totalDuration) / float64(totalRequests)
	}

	errorRate := float64(0)
	if totalRequests > 0 {
		errorRate = float64(totalErrors) / float64(totalRequests) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
		"summary": gin.H{
			"total_requests":     totalRequests,
			"total_errors":       totalErrors,
			"error_rate":         errorRate,
			"avg_response_time":  avgResponseTime,
			"total_endpoints":    len(metrics),
		},
	})
}

// getSlowestEndpoints returns slowest endpoints
func getSlowestEndpoints(c *gin.Context) {
	perfMetricsMux.RLock()
	defer perfMetricsMux.RUnlock()

	var metrics []*PerformanceMetric
	for _, metric := range perfMetrics {
		metrics = append(metrics, metric)
	}

	// Sort by average duration (simple bubble sort for small dataset)
	for i := 0; i < len(metrics)-1; i++ {
		for j := i + 1; j < len(metrics); j++ {
			if metrics[i].AvgDuration < metrics[j].AvgDuration {
				metrics[i], metrics[j] = metrics[j], metrics[i]
			}
		}
	}

	// Get top 10
	limit := 10
	if len(metrics) < limit {
		limit = len(metrics)
	}

	c.JSON(http.StatusOK, gin.H{
		"slowest_endpoints": metrics[:limit],
		"total":             len(metrics),
	})
}

// getMostUsedEndpoints returns most frequently used endpoints
func getMostUsedEndpoints(c *gin.Context) {
	perfMetricsMux.RLock()
	defer perfMetricsMux.RUnlock()

	var metrics []*PerformanceMetric
	for _, metric := range perfMetrics {
		metrics = append(metrics, metric)
	}

	// Sort by count
	for i := 0; i < len(metrics)-1; i++ {
		for j := i + 1; j < len(metrics); j++ {
			if metrics[i].Count < metrics[j].Count {
				metrics[i], metrics[j] = metrics[j], metrics[i]
			}
		}
	}

	// Get top 10
	limit := 10
	if len(metrics) < limit {
		limit = len(metrics)
	}

	c.JSON(http.StatusOK, gin.H{
		"most_used_endpoints": metrics[:limit],
		"total":               len(metrics),
	})
}

// resetPerformanceMetrics resets all performance metrics
func resetPerformanceMetrics(c *gin.Context) {
	perfMetricsMux.Lock()
	defer perfMetricsMux.Unlock()

	perfMetrics = make(map[string]*PerformanceMetric)

	c.JSON(http.StatusOK, gin.H{
		"message": "Performance metrics reset successfully",
	})
}
