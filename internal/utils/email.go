package utils

import (
	"fmt"
	"net/smtp"

	"github.com/capamir/go-api/internal/config"
	"github.com/capamir/go-api/pkg/logger"
)

// EmailService defines what our auth/service layer needs.
type EmailService interface {
	SendVerificationEmail(to, token string) error
}

// ========================================
// Console Email Service (Development)
// ========================================

type ConsoleEmailService struct {
	BaseURL string
}

func NewConsoleEmailService(baseURL string) *ConsoleEmailService {
	return &ConsoleEmailService{BaseURL: baseURL}
}

func (c *ConsoleEmailService) SendVerificationEmail(to, token string) error {
	link := c.BaseURL + "/api/v1/auth/verify?token=" + token
	logger.Info("📧 [DEV] Verification email to %s: %s", to, link)
	return nil
}

// ========================================
// SMTP Email Service (Production)
// ========================================

type SMTPEmailService struct {
	From     string
	Password string
	Host     string
	Port     int
	BaseURL  string
}

func NewSMTPEmailService(cfg *config.Config) *SMTPEmailService {
	return &SMTPEmailService{
		From:     cfg.Email.From,
		Password: cfg.Email.Password,
		Host:     cfg.Email.Host,
		Port:     cfg.Email.Port,
		BaseURL:  cfg.App.URL,
	}
}

func (s *SMTPEmailService) SendVerificationEmail(to, token string) error {
	verificationLink := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", s.BaseURL, token)

	// Email subject and body
	subject := "Verify Your Email"
	body := fmt.Sprintf(`
Hello!

Thank you for registering. Please verify your email by clicking the link below:

%s

If you didn't register, please ignore this email.

Best regards,
E-Commerce API Team
`, verificationLink)

	// Construct email message
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		s.From, to, subject, body,
	))

	// SMTP authentication
	auth := smtp.PlainAuth("", s.From, s.Password, s.Host)

	// Send email
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	if err := smtp.SendMail(addr, auth, s.From, []string{to}, message); err != nil {
		logger.Error("Failed to send email to %s: %v", to, err)
		return err
	}

	logger.Success("✅ Email sent to %s", to)
	return nil
}
