package handler

import (
	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "Registration data"
// @Success 201 {object} service.AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid registration request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Register user
	response, err := h.authService.Register(req)
	if err != nil {
		logger.Error("Registration failed: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondCreated(c, "User registered successfully", response)
}

// Login handles user login
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "Login credentials"
// @Success 200 {object} service.AuthResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid login request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Login user
	response, err := h.authService.Login(req)
	if err != nil {
		logger.Error("Login failed: %v", err)
		utils.RespondUnauthorized(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Login successful", response)
}

// GetProfile returns the authenticated user's profile
// @Summary Get user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Get user profile
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		logger.Error("Failed to get profile for user %d: %v", userID, err)
		utils.RespondNotFound(c, "User not found")
		return
	}

	utils.RespondSuccess(c, "Profile retrieved", user)
}
