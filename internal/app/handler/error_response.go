// Package handler provides standardized error and status response handling for Gin-based APIs.
package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse represents a standardized error response structure.
// It includes:
// - Message: human-readable description of the error
// - Code: machine-readable error code for client-side handling
// - Details: additional context-specific information (e.g., validation errors)
type ErrorResponse struct {
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// StatusResponse represents a standardized success response structure.
// Used for operations that don't return data but need confirmation.
type StatusResponse struct {
	Status string `json:"status"`
}

// ErrorHandler is a struct that encapsulates error handling logic with dependency injection.
// This approach:
// - Promotes testability
// - Allows flexible logging configuration
// - Supports different error handling strategies
type ErrorHandler struct {
	logger *logrus.Logger
}

// NewErrorHandler creates a new ErrorHandler instance with the provided logger.
// Usage:
//
//	errHandler := handler.NewErrorHandler(logrus.StandardLogger())
func NewErrorHandler(logger *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// InternalError handles server-side errors (HTTP 500).
// Parameters:
// - c: Gin context
// - err: original error (for logging purposes)
// - message: user-friendly error message
func (h *ErrorHandler) InternalError(c *gin.Context, err error, message string) {
	// Log with comprehensive context for debugging
	h.logError(c, http.StatusInternalServerError, "internal_error", message, err)

	// Return standardized error response
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
		Message: message,
		Code:    "internal_error",
	})
}

// BadRequest handles client-side validation errors (HTTP 400).
// Parameters:
// - c: Gin context
// - code: machine-readable error code (e.g., "invalid_email")
// - message: user-friendly error description
func (h *ErrorHandler) BadRequest(c *gin.Context, code, message string) {
	h.logError(c, http.StatusBadRequest, code, message, nil)

	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
		Message: message,
		Code:    code,
	})
}

// NotFound handles resource not found errors (HTTP 404).
func (h *ErrorHandler) NotFound(c *gin.Context, message string) {
	h.logError(c, http.StatusNotFound, "not_found", message, nil)

	c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
		Message: message,
		Code:    "not_found",
	})
}

// Unauthorized handles authentication/authorization errors (HTTP 401/403).
func (h *ErrorHandler) Unauthorized(c *gin.Context, code, message string) {
	statusCode := http.StatusUnauthorized
	if code == "forbidden" {
		statusCode = http.StatusForbidden
	}

	h.logError(c, statusCode, code, message, nil)

	c.AbortWithStatusJSON(statusCode, ErrorResponse{
		Message: message,
		Code:    code,
	})
}

// RequestTimeout handles timed out errors (HTTP 408).
func (h *ErrorHandler) RequestTimeout(c *gin.Context, code, message string) {
	h.logError(c, http.StatusRequestTimeout, code, message, nil)

	c.AbortWithStatusJSON(http.StatusRequestTimeout, ErrorResponse{
		Message: message,
		Code:    code,
	})
}

// Success returns a standardized success response with HTTP 200.
func (h *ErrorHandler) Success(c *gin.Context, status string) {
	c.JSON(http.StatusOK, StatusResponse{
		Status: status,
	})
}

// logError records detailed error information with contextual metadata.
// Parameters:
// - c: Gin context (for request metadata)
// - statusCode: HTTP status being returned
// - errorCode: machine-readable error identifier
// - message: user-facing message
// - err: original error (optional)
func (h *ErrorHandler) logError(c *gin.Context, statusCode int, errorCode, message string, err error) {
	fields := logrus.Fields{
		"path":        c.Request.URL.Path,
		"method":      c.Request.Method,
		"status_code": statusCode,
		"error_code":  errorCode,
		"client_ip":   c.ClientIP(),
		"user_agent":  c.Request.UserAgent(),
		"request_id":  c.GetString("request_id"), // If using request ID middleware
	}

	// Include original error details if available
	if err != nil {
		fields["error"] = err.Error()
		fields["stack_trace"] = fmt.Sprintf("%+v", err) // Full stack trace
	}

	h.logger.WithFields(fields).Error(message)
}
