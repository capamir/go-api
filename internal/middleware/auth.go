package middleware

import (
	"strings"

	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header from %s", c.ClientIP())
			utils.RespondUnauthorized(c, "Missing authentication token")
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid Authorization header format from %s", c.ClientIP())
			utils.RespondUnauthorized(c, "Invalid authentication token format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("Invalid token from %s: %v", c.ClientIP(), err)
			utils.RespondUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		logger.Debug("User %d authenticated successfully", claims.UserID)

		c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUserEmail extracts user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	return roleStr, ok
}
