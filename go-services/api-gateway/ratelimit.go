package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter tracks request rates per IP
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

var rateLimiter = &RateLimiter{
	requests: make(map[string][]time.Time),
	limit:    100, // 100 requests
	window:   time.Minute,
}

// checkRateLimit checks if IP has exceeded rate limit
func (rl *RateLimiter) checkRateLimit(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get existing requests for this IP
	requests := rl.requests[ip]

	// Filter out old requests
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if limit exceeded
	if len(validRequests) >= rl.limit {
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests

	return true
}

// getRateLimitInfo returns rate limit info for IP
func (rl *RateLimiter) getRateLimitInfo(ip string) (int, int, time.Time) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	requests := rl.requests[ip]
	var validCount int
	var oldestRequest time.Time

	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validCount++
			if oldestRequest.IsZero() || reqTime.Before(oldestRequest) {
				oldestRequest = reqTime
			}
		}
	}

	remaining := rl.limit - validCount
	if remaining < 0 {
		remaining = 0
	}

	var resetTime time.Time
	if !oldestRequest.IsZero() {
		resetTime = oldestRequest.Add(rl.window)
	}

	return validCount, remaining, resetTime
}

// rateLimitMiddleware enforces rate limiting
func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rateLimiter.checkRateLimit(ip) {
			used, _, resetTime := rateLimiter.getRateLimitInfo(ip)

			c.Header("X-RateLimit-Limit", "100")
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests. Please try again later.",
				"limit":   rateLimiter.limit,
				"used":    used,
				"reset":   resetTime.Format(time.RFC3339),
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		_, remaining, resetTime := rateLimiter.getRateLimitInfo(ip)
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		c.Next()
	}
}

// getRateLimitStatus returns current rate limit status
func getRateLimitStatus(c *gin.Context) {
	ip := c.ClientIP()
	used, remaining, resetTime := rateLimiter.getRateLimitInfo(ip)

	c.JSON(http.StatusOK, gin.H{
		"ip":        ip,
		"limit":     rateLimiter.limit,
		"used":      used,
		"remaining": remaining,
		"reset":     resetTime.Format(time.RFC3339),
		"window":    rateLimiter.window.String(),
	})
}

// Cleanup old rate limit data periodically
func startRateLimitCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			rateLimiter.mu.Lock()
			now := time.Now()
			windowStart := now.Add(-rateLimiter.window)

			for ip, requests := range rateLimiter.requests {
				var validRequests []time.Time
				for _, reqTime := range requests {
					if reqTime.After(windowStart) {
						validRequests = append(validRequests, reqTime)
					}
				}

				if len(validRequests) == 0 {
					delete(rateLimiter.requests, ip)
				} else {
					rateLimiter.requests[ip] = validRequests
				}
			}
			rateLimiter.mu.Unlock()
		}
	}()
}
