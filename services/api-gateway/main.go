package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret       = []byte(getEnv("JWT_SECRET", "your-secret-key-change-in-production"))
	authServiceURL  = getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	threatServiceURL = getEnv("THREAT_SERVICE_URL", "http://localhost:8082")
	networkServiceURL = getEnv("NETWORK_SERVICE_URL", "http://localhost:8083")
	firewallServiceURL = getEnv("FIREWALL_SERVICE_URL", "http://localhost:8084")
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "NetGuard API Gateway v2.0",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "api-gateway",
			"version": "2.0.0",
			"time":    time.Now().Format(time.RFC3339),
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/login", proxyRequest(authServiceURL + "/login"))
	auth.Post("/register", proxyRequest(authServiceURL + "/register"))

	// Protected routes
	protected := api.Group("/", authMiddleware)
	protected.Get("/auth/me", proxyRequest(authServiceURL + "/me"))

	// Dashboard
	protected.Get("/dashboard/stats", getDashboardStats)

	// Alerts
	protected.Get("/alerts", getAlerts)
	protected.Get("/alerts/:id", getAlertById)
	protected.Patch("/alerts/:id/status", updateAlertStatus)

	// Threats
	protected.Get("/threats", proxyRequest(threatServiceURL + "/threats"))
	protected.Get("/threats/:id", proxyRequest(threatServiceURL + "/threats/:id"))
	protected.Post("/threats/:id/block", proxyRequest(threatServiceURL + "/threats/:id/block"))
	protected.Get("/threats/attack-flow", proxyRequest(threatServiceURL + "/threats/attack-flow"))

	// Network
	protected.Get("/network/devices", proxyRequest(networkServiceURL + "/devices"))
	protected.Get("/network/stats", proxyRequest(networkServiceURL + "/stats"))

	// Firewall
	protected.Get("/firewall/rules", proxyRequest(firewallServiceURL + "/rules"))
	protected.Post("/firewall/rules", proxyRequest(firewallServiceURL + "/rules"))
	protected.Put("/firewall/rules/:id", proxyRequest(firewallServiceURL + "/rules/:id"))
	protected.Delete("/firewall/rules/:id", proxyRequest(firewallServiceURL + "/rules/:id"))

	// Users
	protected.Get("/users", getUsers)
	protected.Post("/users", createUser)
	protected.Put("/users/:id", updateUser)
	protected.Delete("/users/:id", deleteUser)

	port := getEnv("PORT", "8080")
	log.Printf("🚀 API Gateway starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - No token provided"})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized - Invalid token"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])
	}

	return c.Next()
}

func proxyRequest(targetURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Replace path params
		url := targetURL
		for _, param := range c.Route().Params {
			url = strings.Replace(url, ":"+param, c.Params(param), 1)
		}

		// Add query params
		if len(c.Queries()) > 0 {
			url += "?" + string(c.Request().URI().QueryString())
		}

		// Create request
		req, err := http.NewRequest(c.Method(), url, bytes.NewReader(c.Body()))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create request"})
		}

		// Copy headers
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// Send request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(503).JSON(fiber.Map{"error": "Service unavailable"})
		}
		defer resp.Body.Close()

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read response"})
		}

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		return c.Status(resp.StatusCode).Send(body)
	}
}

func getDashboardStats(c *fiber.Ctx) error {
	stats := fiber.Map{
		"totalAlerts":   156,
		"activeThreats": 23,
		"blockedAttacks": 1247,
		"systemHealth":  98.5,
		"recentAlerts": []fiber.Map{
			{"id": "ALT-001", "severity": "high", "title": "Suspicious Login Attempt", "time": "2 min ago"},
			{"id": "ALT-002", "severity": "critical", "title": "DDoS Attack Detected", "time": "5 min ago"},
			{"id": "ALT-003", "severity": "medium", "title": "Port Scan Detected", "time": "12 min ago"},
		},
		"threatDistribution": fiber.Map{
			"malware":    45,
			"phishing":   32,
			"ddos":       18,
			"bruteforce": 28,
		},
	}
	return c.JSON(stats)
}

func getAlerts(c *fiber.Ctx) error {
	alerts := []fiber.Map{
		{
			"id": "ALT-001", "title": "Suspicious Login Attempt", "description": "Multiple failed login attempts from 203.0.113.45",
			"severity": "high", "status": "active", "source": "203.0.113.45", "timestamp": time.Now().Add(-2 * time.Minute).Format(time.RFC3339),
		},
		{
			"id": "ALT-002", "title": "DDoS Attack Detected", "description": "High volume of requests from multiple IPs",
			"severity": "critical", "status": "investigating", "source": "Multiple", "timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		},
		{
			"id": "ALT-003", "title": "Port Scan Detected", "description": "Systematic port scanning activity detected",
			"severity": "medium", "status": "active", "source": "198.51.100.23", "timestamp": time.Now().Add(-12 * time.Minute).Format(time.RFC3339),
		},
	}
	return c.JSON(alerts)
}

func getAlertById(c *fiber.Ctx) error {
	id := c.Params("id")
	alert := fiber.Map{
		"id": id, "title": "Suspicious Login Attempt", "description": "Multiple failed login attempts detected",
		"severity": "high", "status": "active", "source": "203.0.113.45", "timestamp": time.Now().Add(-2 * time.Minute).Format(time.RFC3339),
		"details": fiber.Map{"attempts": 15, "lastAttempt": time.Now().Format(time.RFC3339), "targetUser": "admin"},
	}
	return c.JSON(alert)
}

func updateAlertStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var body map[string]string
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	return c.JSON(fiber.Map{"id": id, "status": body["status"], "updated": true})
}

func getUsers(c *fiber.Ctx) error {
	users := []fiber.Map{
		{"id": "USR-001", "email": "admin@netguard.com", "firstName": "Admin", "lastName": "User", "role": "admin", "status": "active", "lastLogin": time.Now().Add(-1 * time.Hour).Format(time.RFC3339)},
		{"id": "USR-002", "email": "analyst@netguard.com", "firstName": "Security", "lastName": "Analyst", "role": "security_analyst", "status": "active", "lastLogin": time.Now().Add(-30 * time.Minute).Format(time.RFC3339)},
	}
	return c.JSON(users)
}

func createUser(c *fiber.Ctx) error {
	var user map[string]interface{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	user["id"] = fmt.Sprintf("USR-%03d", time.Now().Unix()%1000)
	user["status"] = "active"
	user["createdAt"] = time.Now().Format(time.RFC3339)
	return c.Status(201).JSON(user)
}

func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}
	updates["id"] = id
	updates["updatedAt"] = time.Now().Format(time.RFC3339)
	return c.JSON(updates)
}

func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	return c.JSON(fiber.Map{"id": id, "deleted": true})
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error":   err.Error(),
		"code":    code,
		"path":    c.Path(),
		"method":  c.Method(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

# Updated: 2025-08-03T10:00:00

# Updated: 2025-08-03T14:30:00

# Updated: 2025-08-10T09:15:00

# Updated: 2025-08-10T11:00:00

# Updated: 2025-08-11T10:30:00

# Updated: 2025-08-11T14:15:00

# Updated: 2025-08-12T09:45:00

# Updated: 2025-08-12T15:20:00

# Updated: 2025-08-03T10:00:00

# Updated: 2025-08-03T14:30:00

# Updated: 2025-08-10T09:15:00

# Updated: 2025-08-10T11:00:00

# Updated: 2025-08-11T10:30:00

# Updated: 2025-08-11T14:15:00

# Updated: 2025-08-12T09:45:00

# Updated: 2025-08-12T15:20:00

# Updated: 2025-07-03T10:00:00

# Updated: 2025-07-03T14:30:00

# Updated: 2025-07-10T09:15:00

# Updated: 2025-07-10T11:00:00

# Updated: 2025-07-11T10:30:00

# Updated: 2025-07-11T14:15:00

# Updated: 2025-07-12T09:45:00

# Updated: 2025-07-12T15:20:00

# Updated: 2025-07-26T09:15:00

# Updated: 2025-07-27T10:30:00

# Updated: 2025-08-08T14:20:00

# Updated: 2025-08-12T13:45:00

# Updated: 2025-08-14T10:00:00

# Updated: 2025-08-18T15:00:00

# Updated: 2025-08-26T10:00:00

# Updated: 2025-08-29T10:30:00
































