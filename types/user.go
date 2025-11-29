package types

import "time"

// User role 
const (
	UserRoleCustomer = "customer"
	UserRoleAdmin    = "admin"
	UserRoleModerator = "moderator"
)

// User status 
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)

// User represents a user account
type User struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Never expose in JSON
	Phone     string    `json:"phone,omitempty" db:"phone"`
	Role      string    `json:"role" db:"role"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserStore defines the interface for user data operations
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	UpdateUser(User) error
	DeleteUser(id int) error
	GetUsers() ([]*User, error)
}

// RegisterUserPayload represents the data needed to register a new user
type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=8,max=130"` // ✅ Fixed: min=8
	Phone     string `json:"phone" validate:"omitempty,e164"` // E.164 format (+1234567890)
}

// LoginUserPayload represents the data needed to login
type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserPayload represents the data that can be updated for a user
type UpdateUserPayload struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     *string `json:"email" validate:"omitempty,email,max=100"`
	Phone     *string `json:"phone" validate:"omitempty,e164"`
}

// ChangePasswordPayload represents the data needed to change password
type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=130"`
}

// UserProfileResponse represents user data returned to the client (safe)
type UserProfileResponse struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	Token string              `json:"token"`
	User  UserProfileResponse `json:"user"`
}
