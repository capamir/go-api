package user

import (
	"fmt"
	"net/http"

	"github.com/capamir/go-api/configs"
	"github.com/capamir/go-api/services/auth"
	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Public routes
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)

	// Protected routes (require JWT)
	router.HandleFunc("/users/me", auth.WithJWTAuth(h.handleGetCurrentUser, h.store)).Methods(http.MethodGet)
}

// handleLogin authenticates a user and returns a JWT token
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse request payload
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.S.Warnf("Invalid JSON in login request: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.S.Warnf("Login validation failed: %v", validationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrors))
		return
	}

	// Get user by email
	user, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		// Don't reveal whether email exists or not (security)
		utils.S.Warnf("Login attempt failed for email: %s", payload.Email)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	// Compare passwords
	if !auth.ComparePasswords(user.Password, []byte(payload.Password)) {
		utils.S.Warnf("Invalid password attempt for user %d (%s)", user.ID, user.Email)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid email or password"))
		return
	}

	// Generate JWT token
	token, err := auth.CreateJWT([]byte(configs.Envs.JWTSecret), user.ID)
	if err != nil {
		utils.S.Errorf("Failed to create JWT for user %d: %v", user.ID, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to generate token"))
		return
	}

	utils.S.Successf("User %d (%s) logged in successfully", user.ID, user.Email)

	// Return token and safe user data
	response := types.AuthResponse{
		Token: token,
		User: types.UserProfileResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// handleRegister creates a new user account
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// Parse request payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.S.Warnf("Invalid JSON in register request: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.S.Warnf("Registration validation failed: %v", validationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrors))
		return
	}

	// Check if user already exists
	existingUser, err := h.store.GetUserByEmail(payload.Email)
	if err == nil && existingUser != nil {
		// User exists - don't reveal this for security (email enumeration)
		utils.S.Warnf("Registration attempt with existing email: %s", payload.Email)
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("email already registered"))
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.S.Errorf("Failed to hash password during registration: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to process registration"))
		return
	}

	// Create new user
	user := types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}

	if err := h.store.CreateUser(user); err != nil {
		utils.S.Errorf("Failed to create user with email %s: %v", payload.Email, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create user"))
		return
	}

	utils.S.Successf("New user registered with email: %s", payload.Email)

	// Return success message
	utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
		"email":   payload.Email,
	})
}

// handleGetCurrentUser returns the authenticated user's profile
func (h *Handler) handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.S.Errorf("Failed to get user ID from context: %v", err)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
		return
	}

	// Get user from database
	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.S.Errorf("Failed to get user %d: %v", userID, err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user not found"))
		return
	}

	// Return safe user profile (without password)
	profile := types.UserProfileResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	utils.WriteJSON(w, http.StatusOK, profile)
}
