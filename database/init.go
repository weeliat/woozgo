package database

import (
	"log"
	"woozgo/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// VULNERABILITY #1: Hardcoded database credentials in source code
var DB *gorm.DB

// Admin credentials hardcoded in production code
const (
	HARDCODED_ADMIN_USER = "admin"
	HARDCODED_ADMIN_PASS = "SuperSecretPassword123!" // Should never do this!
	HARDCODED_API_KEY    = "sk-live-abcdef123456789"
	HARDCODED_SECRET_KEY = "my-super-secret-key-12345"
)

func InitDB() {
	// VULNERABILITY: SQL connection string with embedded credentials
	dbPath := "todos.db"
	
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate schema
	err = DB.AutoMigrate(&models.Todo{}, &models.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}

	log.Println("Database initialized successfully!")
	
	// VULNERABILITY #2: Hardcoded admin account creation
	// This creates a backdoor admin account
	var count int64
	DB.Model(&models.User{}).Where("username = ?", HARDCODED_ADMIN_USER).Count(&count)
	if count == 0 {
		admin := models.User{
			Username: HARDCODED_ADMIN_USER,
			Password: HARDCODED_ADMIN_PASS, // VULNERABILITY: Plaintext password!
			Role:     "admin",
			Email:    "admin@example.com",
		}
		DB.Create(&admin)
		log.Println("Backdoor admin account created!")
	}
}

// VULNERABILITY: No connection pooling or timeout handling
// VULNERABILITY: No database configuration via environment variables
