package utils

import "github.com/capamir/go-api/pkg/logger"

// EmailService defines what our auth/service layer needs.
type EmailService interface {
	SendVerificationEmail(to, token string) error
	// Later: SendPasswordResetEmail(to, token string) error
}

// ConsoleEmailService is a development implementation that just logs.
type ConsoleEmailService struct {
	BaseURL string // e.g. http://localhost:8080
}

// NewConsoleEmailService creates a new console email service.
func NewConsoleEmailService(baseURL string) *ConsoleEmailService {
	return &ConsoleEmailService{BaseURL: baseURL}
}

// SendVerificationEmail logs the verification link instead of sending an email.
func (c *ConsoleEmailService) SendVerificationEmail(to, token string) error {
	link := c.BaseURL + "/api/v1/auth/verify?token=" + token
	logger.Info("📧 [DEV] Verification email to %s: %s", to, link)
	return nil
}
