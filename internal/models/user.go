package models

import (
	"time"
	"gorm.io/gorm"
)

// User status constants
const (
	UserStatusPending   = "pending"   // Email not verified yet
	UserStatusActive    = "active"    // Email verified, can login
	UserStatusSuspended = "suspended" // Temporarily blocked by admin
	UserStatusBanned    = "banned"    // Permanently blocked
	UserStatusInactive  = "inactive"  // User deactivated their own account
)

// User role constants
const (
	UserRoleCustomer  = "customer"
	UserRoleAdmin     = "admin"
	UserRoleModerator = "moderator"
)

// User represents a user in the system
type User struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	FirstName         string         `gorm:"size:100;not null" json:"first_name"`
	LastName          string         `gorm:"size:100;not null" json:"last_name"`
	Email             string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	EmailVerified     bool           `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt   *time.Time     `json:"email_verified_at,omitempty"`
	VerificationToken string         `gorm:"size:255;index" json:"-"` // ← Added index for faster lookup
	Password          string         `gorm:"size:255;not null" json:"-"`
	Phone             string         `gorm:"size:20" json:"phone,omitempty"`
	Role              string         `gorm:"size:20;default:'customer';index" json:"role"` // ← Added index
	Status            string         `gorm:"size:20;default:'pending';index" json:"status"` // ← Added index
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Set default values if not provided
	if u.Role == "" {
		u.Role = UserRoleCustomer
	}
	if u.Status == "" {
		u.Status = UserStatusPending
	}
	return nil
}

// IsEmailVerified checks if user verified their email
func (u *User) IsEmailVerified() bool {
	return u.EmailVerified
}

// IsActive checks if user can use the system
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive && u.EmailVerified
}

// IsSuspended checks if user is suspended
func (u *User) IsSuspended() bool {
	return u.Status == UserStatusSuspended
}

// IsBanned checks if user is banned
func (u *User) IsBanned() bool {
	return u.Status == UserStatusBanned
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// CanLogin checks if user can login (active + email verified + not suspended/banned)
func (u *User) CanLogin() bool {
	if u.Status == UserStatusBanned || u.Status == UserStatusSuspended {
		return false
	}
	return u.EmailVerified && u.Status == UserStatusActive
}

// UserResponse is a safe representation of User for API responses
type UserResponse struct {
	ID              uint       `json:"id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	Email           string     `json:"email"`
	EmailVerified   bool       `json:"email_verified"` // ← Added
	Phone           string     `json:"phone,omitempty"`
	Role            string     `json:"role"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"` // ← Added
}

// ToResponse converts User to UserResponse (removes sensitive data)
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:              u.ID,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Email:           u.Email,
		EmailVerified:   u.EmailVerified,
		Phone:           u.Phone,
		Role:            u.Role,
		Status:          u.Status,
		CreatedAt:       u.CreatedAt,
		EmailVerifiedAt: u.EmailVerifiedAt,
	}
}

// FullName returns user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
