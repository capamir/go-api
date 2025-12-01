package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSecureToken generates a cryptographically secure random token
// length is in bytes (output will be length*2 hex characters).
func GenerateSecureToken(length int) (string, error) {
	if length <= 0 {
		length = 32
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken returns a 64-char hex token (32 random bytes).
func GenerateVerificationToken() (string, error) {
	return GenerateSecureToken(32)
}
