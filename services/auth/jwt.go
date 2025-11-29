package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/capamir/go-api/configs"
	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
	"github.com/golang-jwt/jwt/v5"
)

// Context key for storing user information
type contextKey string
const UserKey contextKey = "userID"

// CustomClaims defines the structure of JWT claims
type CustomClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// CreateJWT generates a new JWT token for the given user ID
func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Duration(configs.Envs.JWTExpirationInSeconds) * time.Second

	// Use structured claims with standard JWT fields
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-api",                    // Your API name
			Subject:   fmt.Sprintf("%d", userID),   // User ID as subject
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateJWT validates and parses a JWT token string
func ValidateJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.Envs.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims with type safety
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	// Additional validation (automatically handled by jwt library, but explicit for clarity)
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// WithJWTAuth is a middleware that validates JWT tokens and adds user info to context
func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from request
		tokenString := utils.GetTokenFromRequest(r)
		if tokenString == "" {
			utils.S.Warnf("Missing authentication token from %s", r.RemoteAddr)
			writeUnauthorized(w, fmt.Errorf("missing authentication token"))
			return
		}

		// Validate token
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			utils.S.Warnf("Invalid token from %s: %v", r.RemoteAddr, err)
			writeUnauthorized(w, fmt.Errorf("invalid or expired token"))
			return
		}

		// Verify user exists in database
		user, err := store.GetUserByID(claims.UserID)
		if err != nil {
			utils.S.Errorf("User %d from token not found in database: %v", claims.UserID, err)
			writeUnauthorized(w, fmt.Errorf("user not found"))
			return
		}

		// Optional: Check if user is active/not banned
		// if user.Status == "banned" {
		//     writeUnauthorized(w, fmt.Errorf("user account is disabled"))
		//     return
		// }

		utils.S.Debugf("Authenticated user %d (%s) for %s %s", user.ID, user.Email, r.Method, r.URL.Path)

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserKey, user.ID)
		r = r.WithContext(ctx)

		// Call the protected handler
		handlerFunc(w, r)
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// writeUnauthorized sends a 401 Unauthorized response
func writeUnauthorized(w http.ResponseWriter, err error) {
	utils.WriteError(w, http.StatusUnauthorized, err)
}

// RefreshToken generates a new token for an existing valid token (token refresh flow)
func RefreshToken(tokenString string) (string, error) {
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("cannot refresh invalid token: %w", err)
	}

	// Generate new token with same user ID
	return CreateJWT([]byte(configs.Envs.JWTSecret), claims.UserID)
}
