package service

import (
	"errors"
	"fmt"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

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
	Token string               `json:"token"`
	User  models.UserResponse `json:"user"`
}

// Register creates a new user account
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

	// Create user
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
		Phone:     req.Phone,
		Role:      "customer",
		Status:    "active",
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user: %v", err)
		return nil, errors.New("failed to register user")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		return nil, errors.New("failed to generate authentication token")
	}

	logger.Success("User registered successfully: %s", user.Email)

	return &AuthResponse{
		Token: token,
		User:  user.ToResponse(),
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		logger.Error("Database error during login: %v", err)
		return nil, errors.New("login failed")
	}
	if user == nil {
		logger.Warn("Login attempt with non-existent email: %s", req.Email)
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if user.Status != "active" {
		logger.Warn("Login attempt for inactive user: %s", req.Email)
		return nil, errors.New("account is not active")
	}

	// Verify password
	if !utils.ComparePassword(user.Password, req.Password) {
		logger.Warn("Invalid password attempt for user: %s", req.Email)
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
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

// GetUserByID retrieves a user by ID
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
