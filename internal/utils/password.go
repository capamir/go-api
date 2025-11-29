package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	// MinCost is the minimum bcrypt cost (for testing only)
	MinCost = bcrypt.MinCost // 4
	// DefaultCost is the default bcrypt cost (recommended)
	DefaultCost = bcrypt.DefaultCost // 10
	// MaxCost is expensive but very secure
	MaxCost = 14
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Validate password
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	// Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plain password
func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// ValidatePassword checks if password meets requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return errors.New("password must be less than 72 characters") // bcrypt limit
	}
	return nil
}
