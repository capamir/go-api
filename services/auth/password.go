package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinCost     = bcrypt.MinCost     // 4 (for testing only)
	DefaultCost = bcrypt.DefaultCost // 10 (recommended)
	MaxCost     = 14                 // Very secure but slower
)

// ValidatePassword checks if a password meets minimum security requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	
	if len(password) > 72 {
		return errors.New("password must be less than 72 characters") // bcrypt limit
	}

	return nil
}

// HashPassword validates and hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	// Validate before hashing
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}
