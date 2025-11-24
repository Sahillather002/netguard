package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Company      string    `json:"company"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// Session represents an active user session
type Session struct {
	Token        string
	RefreshToken string
	UserID       string
	ExpiresAt    time.Time
}

// In-memory storage (replace with real database in production)
var (
	users    = make(map[string]*User)
	sessions = make(map[string]*Session)
	usersMux sync.RWMutex
)

func init() {
	// Create a default test user
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	users["test@example.com"] = &User{
		ID:           "user_001",
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: string(hash),
		Company:      "Test Company",
		Role:         "admin",
		CreatedAt:    time.Now(),
	}
}

// generateToken generates a random token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// verifyPassword verifies a password against a hash
func verifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// authMiddleware validates JWT token
// authOrAPIKeyMiddleware accepts either JWT token or API key
func authOrAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try API key first
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			// Validate API key
			apiKeysMux.RLock()
			key, exists := apiKeys[apiKey]
			apiKeysMux.RUnlock()

			if exists && key.Enabled && time.Now().Before(key.ExpiresAt) {
				// API key is valid, set user context
				usersMux.RLock()
				user, userExists := users[key.UserID]
				usersMux.RUnlock()

				if userExists {
					c.Set("user", user)
					c.Set("user_id", user.ID)
					c.Set("auth_method", "api_key")
					c.Next()
					return
				}
			}
		}

		// Try JWT token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]

				// Validate token
				usersMux.RLock()
				session, exists := sessions[token]
				usersMux.RUnlock()

				if exists && time.Now().Before(session.ExpiresAt) {
					// Get user
					usersMux.RLock()
					user, userExists := users[session.UserID]
					usersMux.RUnlock()

					if userExists {
						c.Set("user", user)
						c.Set("user_id", user.ID)
						c.Set("auth_method", "jwt")
						c.Next()
						return
					}
				}
			}
		}

		// Neither API key nor JWT token is valid
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required. Provide valid JWT token or API key"})
		c.Abort()
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		usersMux.RLock()
		session, exists := sessions[token]
		usersMux.RUnlock()

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if token is expired
		if time.Now().After(session.ExpiresAt) {
			usersMux.Lock()
			delete(sessions, token)
			usersMux.Unlock()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Get user
		usersMux.RLock()
		user, exists := users[session.UserID]
		usersMux.RUnlock()

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Next()
	}
}

// login handles user login
func login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Find user
	usersMux.RLock()
	user, exists := users[req.Email]
	usersMux.RUnlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if !verifyPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate tokens
	token := generateToken()
	refreshToken := generateToken()

	// Create session
	session := &Session{
		Token:        token,
		RefreshToken: refreshToken,
		UserID:       user.Email,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	usersMux.Lock()
	sessions[token] = session
	sessions[refreshToken] = session
	usersMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"expires_in":    86400,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"company": user.Company,
			"role":    user.Role,
		},
	})
}

// register handles user registration
func register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
		Name     string `json:"name" binding:"required"`
		Company  string `json:"company"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Check if user already exists
	usersMux.RLock()
	_, exists := users[req.Email]
	usersMux.RUnlock()

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hash, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user := &User{
		ID:           fmt.Sprintf("user_%d", time.Now().Unix()),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hash,
		Company:      req.Company,
		Role:         req.Role,
		CreatedAt:    time.Now(),
	}

	usersMux.Lock()
	users[req.Email] = user
	usersMux.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
		"user": gin.H{
			"id":      user.ID,
			"email":   user.Email,
			"name":    user.Name,
			"company": user.Company,
			"role":    user.Role,
		},
	})
}

// refreshToken handles token refresh
func refreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate refresh token
	usersMux.RLock()
	session, exists := sessions[req.RefreshToken]
	usersMux.RUnlock()

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access token
	newToken := generateToken()
	newSession := &Session{
		Token:        newToken,
		RefreshToken: session.RefreshToken,
		UserID:       session.UserID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	usersMux.Lock()
	sessions[newToken] = newSession
	usersMux.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"token":      newToken,
		"expires_in": 86400,
	})
}

// logout handles user logout
func logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 {
			token := parts[1]
			usersMux.Lock()
			delete(sessions, token)
			usersMux.Unlock()
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
