package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

var (
	db        *sql.DB
	jwtSecret = []byte(getEnv("JWT_SECRET", "your-secret-key-change-in-production"))
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Role      string    `json:"role"`
	Company   string    `json:"company"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Company   string `json:"company"`
	Role      string `json:"role"`
}

func main() {
	// Initialize database
	initDB()
	defer db.Close()

	app := fiber.New(fiber.Config{
		AppName: "NetGuard Auth Service v2.0",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy", "service": "auth-service"})
	})

	app.Post("/login", login)
	app.Post("/register", register)
	app.Get("/me", getCurrentUser)

	port := getEnv("PORT", "8081")
	log.Printf("🔐 Auth Service starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func initDB() {
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		log.Println("⚠️  No DATABASE_URL set, using in-memory mode")
		return
	}

	var err error
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("⚠️  Database connection failed: %v (using in-memory mode)", err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Printf("⚠️  Database ping failed: %v (using in-memory mode)", err)
		db = nil
		return
	}

	log.Println("✓ Database connected successfully")
	createTables()
}

func createTables() {
	if db == nil {
		return
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		role VARCHAR(50) NOT NULL,
		company VARCHAR(255),
		status VARCHAR(20) DEFAULT 'active',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		log.Printf("Error creating tables: %v", err)
	} else {
		log.Println("✓ Database tables ready")
	}
}

func login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// For demo, accept admin@netguard.com / password
	if req.Email == "admin@netguard.com" && req.Password == "password" {
		user := User{
			ID:        "user-admin-001",
			Email:     req.Email,
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
			Company:   "NetGuard",
			Status:    "active",
		}

		token, err := generateToken(user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		return c.JSON(fiber.Map{
			"token": token,
			"user":  user,
		})
	}

	// Try database if available
	if db != nil {
		var user User
		var passwordHash string

		err := db.QueryRow(`
			SELECT id, email, password_hash, first_name, last_name, role, company, status
			FROM users WHERE email = $1 AND status = 'active'
		`, req.Email).Scan(&user.ID, &user.Email, &passwordHash, &user.FirstName, &user.LastName, &user.Role, &user.Company, &user.Status)

		if err == nil {
			if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err == nil {
				// Update last login
				db.Exec("UPDATE users SET last_login = $1 WHERE id = $2", time.Now(), user.ID)

				token, err := generateToken(user)
				if err != nil {
					return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
				}

				return c.JSON(fiber.Map{
					"token": token,
					"user":  user,
				})
			}
		}
	}

	return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
}

func register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	user := User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		Company:   req.Company,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	// Save to database if available
	if db != nil {
		_, err := db.Exec(`
			INSERT INTO users (id, email, password_hash, first_name, last_name, role, company, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, user.ID, user.Email, string(hashedPassword), user.FirstName, user.LastName, user.Role, user.Company, user.Status, user.CreatedAt)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Email already exists or database error"})
		}
	}

	token, err := generateToken(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(201).JSON(fiber.Map{
		"token": token,
		"user":  user,
	})
}

func getCurrentUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tokenString := authHeader[7:] // Remove "Bearer "
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	user := User{
		ID:        claims["user_id"].(string),
		Email:     claims["email"].(string),
		FirstName: claims["first_name"].(string),
		LastName:  claims["last_name"].(string),
		Role:      claims["role"].(string),
		Company:   claims["company"].(string),
	}

	return c.JSON(user)
}

func generateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"company":    user.Company,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"iat":        time.Now().Unix(),
	})

	return token.SignedString(jwtSecret)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

# Updated: 2025-08-04T09:45:00

# Updated: 2025-08-04T11:20:00

# Updated: 2025-08-04T15:00:00

# Updated: 2025-08-04T09:45:00

# Updated: 2025-08-04T11:20:00

# Updated: 2025-08-04T15:00:00

# Updated: 2025-07-04T09:45:00

# Updated: 2025-07-04T11:20:00

# Updated: 2025-07-04T15:00:00

# Updated: 2025-08-08T09:45:00

# Updated: 2025-08-12T13:45:00

# Updated: 2025-08-16T11:00:00

# Updated: 2025-08-26T15:30:00

# Updated: 2025-08-29T10:30:00

# Updated: 2025-09-12T10:30:00

# Updated: 2025-10-01T10:15:00




















