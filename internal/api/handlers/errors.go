package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Common API errors
var (
	ErrBadRequest          = NewAPIError(fiber.StatusBadRequest, "BAD_REQUEST", "Invalid request parameters")
	ErrNotFound            = NewAPIError(fiber.StatusNotFound, "NOT_FOUND", "Resource not found")
	ErrInternalServer      = NewAPIError(fiber.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	ErrTooManyRequests     = NewAPIError(fiber.StatusTooManyRequests, "RATE_LIMITED", "Too many requests")
	ErrValidationFailed    = NewAPIError(fiber.StatusUnprocessableEntity, "VALIDATION_FAILED", "Validation failed")
	ErrServiceUnavailable  = NewAPIError(fiber.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service temporarily unavailable")
)

// APIError represents a structured API error
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, code, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// WithDetails adds details to an error
func (e *APIError) WithDetails(details string) *APIError {
	return &APIError{
		StatusCode: e.StatusCode,
		Code:       e.Code,
		Message:    e.Message,
		Details:    details,
	}
}

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

// SendError sends an error response
func SendError(c *fiber.Ctx, err *APIError) error {
	return c.Status(err.StatusCode).JSON(ErrorResponse{Error: err})
}

// SendErrorWithDetails sends an error with additional details
func SendErrorWithDetails(c *fiber.Ctx, err *APIError, details string) error {
	return c.Status(err.StatusCode).JSON(ErrorResponse{
		Error: err.WithDetails(details),
	})
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors holds multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// SendValidationError sends a validation error response
func SendValidationError(c *fiber.Ctx, errors []ValidationError) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "VALIDATION_FAILED",
			"message": "One or more fields failed validation",
			"errors":  errors,
		},
	})
}
