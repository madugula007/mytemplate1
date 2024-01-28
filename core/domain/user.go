package domain

import (
	"time"
)

// UserRole is an enum for user's role
type UserRole string

// UserRole enum values
const (
	Admin   UserRole = "admin"
	Cashier UserRole = "cashier"
)

// User is an entity that represents a user
type User struct {
	ID        uint64    `json:"id" db:"id" select:"id" `
	Name      string    `json:"name" insert:"name" select:"name" insert_pickup:"name"`
	Email     string    `json:"email" insert:"email" select:"email"`
	Password  string    `json:"password" insert:"password" select:"password"`
	Role      UserRole  `json:"role" db:"role" select:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at" select:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" select:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=5" u:"N1" db:"name" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" db:"email" example:"test@example.com"`
	Password string `json:"password" validate:"required,min=8" u:"P1" db:"password" example:"12345678"`
	Check    int    `json:"check"  validate:"required,myvalidate"`
}
