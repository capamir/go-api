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

// ========================================
// EXISTING HANDLERS
// ========================================

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

	// Updated message to reflect pending verification
	utils.RespondCreated(c, "Registration successful! Please check your email to verify your account.", response)
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

// ========================================
// 🆕 EMAIL VERIFICATION HANDLERS
// ========================================

// VerifyEmail handles email verification via token
// @Summary Verify user email
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} service.VerificationResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /auth/verify [get]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	
	// Validate token parameter
	if token == "" {
		logger.Warn("Email verification attempted without token")
		utils.RespondBadRequest(c, "Verification token is required")
		return
	}

	// Verify email
	response, err := h.authService.VerifyEmail(token)
	if err != nil {
		logger.Error("Email verification failed: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	logger.Success("Email verified successfully for user: %s", response.User.Email)
	utils.RespondSuccess(c, response.Message, response)
}

// ResendVerification handles resending verification email
// @Summary Resend verification email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendVerificationRequest true "Email address"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req ResendVerificationRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid resend verification request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Resend verification email
	if err := h.authService.ResendVerificationEmail(req.Email); err != nil {
		logger.Error("Failed to resend verification email: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// Generic success message (don't reveal if email exists)
	utils.RespondSuccess(c, "If the email exists and is not verified, a verification link has been sent.", nil)
}

// ========================================
// REQUEST DTOs (for this handler)
// ========================================

// ResendVerificationRequest represents the request body for resending verification
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}
