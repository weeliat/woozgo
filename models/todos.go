package models

import (
	"time"
)

// Todo represents a task item
type Todo struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Description string   `json:"description"`
	Completed bool      `gorm:"default:false" json:"completed"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// User represents a user account
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `json:"password"` // VULNERABILITY: Weak hashing in handlers
	Email     string    `json:"email"`
	Role      string    `gorm:"default:user" json:"role"` // user or admin
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthToken represents a session token
type AuthToken struct {
	UserID   uint      `json:"user_id"`
	Token    string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VULNERABILITY: No input validation on model fields
// VULNERABILITY: No length limits on strings
