package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RespondJSON sends a JSON response
func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// RespondSuccess sends a success response
func RespondSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondError sends an error response
func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// RespondCreated sends a 201 Created response
func RespondCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondBadRequest sends a 400 Bad Request response
func RespondBadRequest(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, message)
}

// RespondUnauthorized sends a 401 Unauthorized response
func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message)
}

// RespondForbidden sends a 403 Forbidden response
func RespondForbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, message)
}

// RespondNotFound sends a 404 Not Found response
func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message)
}

// RespondInternalError sends a 500 Internal Server Error response
func RespondInternalError(c *gin.Context, message string) {
	RespondError(c, http.StatusInternalServerError, message)
}
