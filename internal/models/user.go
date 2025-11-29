package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	FirstName string         `gorm:"size:100;not null" json:"first_name"`
	LastName  string         `gorm:"size:100;not null" json:"last_name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"` // Never send password in JSON
	Phone     string         `gorm:"size:20" json:"phone,omitempty"`
	Role      string         `gorm:"size:20;default:'customer'" json:"role"`
	Status    string         `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Set default values if not provided
	if u.Role == "" {
		u.Role = "customer"
	}
	if u.Status == "" {
		u.Status = "active"
	}
	return nil
}

// UserResponse is a safe representation of User for API responses
type UserResponse struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse (removes sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
