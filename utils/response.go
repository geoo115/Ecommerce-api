package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// SendSuccess sends a successful response
func SendSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Code:    status,
	})
}

// SendError sends an error response
func SendError(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error:   message,
		Code:    status,
	})
}

// SendValidationError sends a validation error response
func SendValidationError(c *gin.Context, message string) {
	SendError(c, http.StatusBadRequest, message)
}

// SendUnauthorized sends an unauthorized error response
func SendUnauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized access"
	}
	SendError(c, http.StatusUnauthorized, message)
}

// SendNotFound sends a not found error response
func SendNotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	SendError(c, http.StatusNotFound, message)
}

// SendInternalError sends an internal server error response
func SendInternalError(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	SendError(c, http.StatusInternalServerError, message)
}

// SendConflict sends a conflict error response
func SendConflict(c *gin.Context, message string) {
	if message == "" {
		message = "Resource conflict"
	}
	SendError(c, http.StatusConflict, message)
}

// SendForbidden sends a forbidden error response
func SendForbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Access forbidden"
	}
	SendError(c, http.StatusForbidden, message)
}
