package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo     *repository.UserRepository
	emailService utils.EmailService
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, emailService utils.EmailService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

// ========================================
// DTOs (Data Transfer Objects)
// ========================================

// RegisterRequest represents registration data
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string `json:"last_name" binding:"required,min=2,max=100"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=72"`
	Phone     string `json:"phone" binding:"omitempty"`
}

// LoginRequest represents login data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string              `json:"token,omitempty"` // Optional for pending users
	User  models.UserResponse `json:"user"`
}

// VerificationResponse represents email verification response
type VerificationResponse struct {
	Message string              `json:"message"`
	User    models.UserResponse `json:"user"`
}

// ========================================
// REGISTER (with Email Verification)
// ========================================

// Register creates a new user account (pending until email verified)
func (s *AuthService) Register(req RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		logger.Error("Failed to check email existence: %v", err)
		return nil, errors.New("failed to register user")
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password: %v", err)
		return nil, errors.New("failed to register user")
	}

	// Generate verification token
	verificationToken, err := utils.GenerateVerificationToken()
	if err != nil {
		logger.Error("Failed to generate verification token: %v", err)
		return nil, errors.New("failed to register user")
	}

	// Create user with pending status
	user := &models.User{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		Password:          hashedPassword,
		Phone:             req.Phone,
		Role:              models.UserRoleCustomer,
		Status:            models.UserStatusPending,
		EmailVerified:     false,
		VerificationToken: verificationToken,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user: %v", err)
		return nil, errors.New("failed to register user")
	}

	// Send verification email via EmailService
	if s.emailService != nil {
		if err := s.emailService.SendVerificationEmail(user.Email, verificationToken); err != nil {
			logger.Error("Failed to send verification email to %s: %v", user.Email, err)
			// Do NOT fail the registration because of email issues
		}
	} else {
		// Fallback: log if no email service provided
		logger.Warn("EmailService is nil, cannot send verification email for %s", user.Email)
	}

	logger.Success("User registered (pending verification): %s", user.Email)

	// Don't generate JWT token yet - user must verify email first
	return &AuthResponse{
		User: user.ToResponse(),
	}, nil
}

// ========================================
// LOGIN (Check Email Verification)
// ========================================

func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		logger.Error("Database error during login: %v", err)
		return nil, errors.New("login failed")
	}
	if user == nil {
		logger.Warn("Login attempt with non-existent email: %s", req.Email)
		return nil, errors.New("invalid email or password")
	}

	if !user.EmailVerified {
		logger.Warn("Login attempt with unverified email: %s", req.Email)
		return nil, errors.New("please verify your email before logging in")
	}

	if !user.CanLogin() {
		logger.Warn("Login attempt for inactive/banned user: %s (status: %s)", req.Email, user.Status)
		return nil, errors.New("account is not active")
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		logger.Warn("Invalid password attempt for user: %s", req.Email)
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error("Failed to generate token for user %s: %v", user.Email, err)
		return nil, errors.New("failed to generate authentication token")
	}

	logger.Success("User logged in successfully: %s", user.Email)

	return &AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// ========================================
// EMAIL VERIFICATION
// ========================================

func (s *AuthService) VerifyEmail(token string) (*VerificationResponse, error) {
	user, err := s.userRepo.GetByVerificationToken(token)
	if err != nil {
		logger.Error("Database error during email verification: %v", err)
		return nil, errors.New("verification failed")
	}
	if user == nil {
		logger.Warn("Invalid verification token attempted: %s", token)
		return nil, errors.New("invalid or expired verification token")
	}

	if user.EmailVerified {
		logger.Info("User %s already verified", user.Email)
		return &VerificationResponse{
			Message: "Email already verified",
			User:    user.ToResponse(),
		}, nil
	}

	now := time.Now()
	user.EmailVerified = true
	user.EmailVerifiedAt = &now
	user.Status = models.UserStatusActive
	user.VerificationToken = ""

	if err := s.userRepo.Update(user); err != nil {
		logger.Error("Failed to update user verification status: %v", err)
		return nil, errors.New("verification failed")
	}

	logger.Success("✅ Email verified for user: %s", user.Email)

	return &VerificationResponse{
		Message: "Email verified successfully! You can now login.",
		User:    user.ToResponse(),
	}, nil
}

// ========================================
// RESEND VERIFICATION EMAIL
// ========================================

func (s *AuthService) ResendVerificationEmail(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		logger.Error("Database error during resend verification: %v", err)
		return errors.New("failed to resend verification email")
	}
	if user == nil {
		logger.Warn("Resend verification attempted for non-existent email: %s", email)
		return errors.New("if the email exists, a verification link has been sent")
	}

	if user.EmailVerified {
		logger.Info("Resend verification attempted for already verified user: %s", email)
		return errors.New("email is already verified")
	}

	if time.Since(user.CreatedAt) < 1*time.Minute {
		logger.Warn("Resend verification attempted too soon for: %s", email)
		return errors.New("please wait before requesting a new verification email")
	}

	newToken, err := utils.GenerateVerificationToken()
	if err != nil {
		logger.Error("Failed to generate new verification token: %v", err)
		return errors.New("failed to resend verification email")
	}

	user.VerificationToken = newToken
	if err := s.userRepo.Update(user); err != nil {
		logger.Error("Failed to update verification token: %v", err)
		return errors.New("failed to resend verification email")
	}

	if s.emailService != nil {
		if err := s.emailService.SendVerificationEmail(user.Email, newToken); err != nil {
			logger.Error("Failed to resend verification email to %s: %v", user.Email, err)
			return errors.New("failed to resend verification email")
		}
	} else {
		logger.Warn("EmailService is nil, cannot resend verification email for %s", user.Email)
	}

	logger.Success("Verification email resent to: %s", user.Email)
	return nil
}

// ========================================
// HELPER METHODS
// ========================================

func (s *AuthService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get user %d: %v", id, err)
		return nil, fmt.Errorf("failed to get user")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}
