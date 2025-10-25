// Core Application Domain-Agnostic errors
package errors

import (
	"fmt"
	"net/http"
)

// AppError is the base application error
type AppError struct {
	Code    string
	Message string
	Cause   error
	Status  int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// New creates a new AppError
func New(code, message string, status int, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Status:  status,
	}
}

// Common Application Errors
var (
	// 400 Bad Request
	ErrInvalidInput = New("ERR_INVALID_INPUT", "Invalid input data", http.StatusBadRequest, nil)
	ErrValidation   = New("ERR_VALIDATION", "Validation failed", http.StatusBadRequest, nil)

	// 401 Unauthorized
	ErrUnauthorized = New("ERR_UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized, nil)
	ErrInvalidToken = New("ERR_INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized, nil)

	// 404 Not Found
	ErrNotFound = New("ERR_NOT_FOUND", "Resource not found", http.StatusNotFound, nil)

	// 409 Conflict
	ErrAlreadyExists = New("ERR_ALREADY_EXISTS", "Resource already exists", http.StatusConflict, nil)
	ErrConflict      = New("ERR_CONFLICT", "Operation conflicts with current state", http.StatusConflict, nil)

	// 413 Payload Too Large
	ErrPayloadTooLarge = New("ERR_PAYLOAD_TOO_LARGE", "Payload too large", http.StatusRequestEntityTooLarge, nil)

	// 422 Unprocessable Entity
	ErrUnprocessable = New("ERR_UNPROCESSABLE", "Unable to process request", http.StatusUnprocessableEntity, nil)

	// 429 Too Many Requests
	ErrRateLimited = New("ERR_RATE_LIMITED", "Too many requests", http.StatusTooManyRequests, nil)

	// 500 Internal Server Error
	ErrInternal = New("ERR_INTERNAL", "Internal server error", http.StatusInternalServerError, nil)
	ErrDatabase = New("ERR_DATABASE", "Database error", http.StatusInternalServerError, nil)

	// 503 Service Unavailable
	ErrServiceUnavailable = New("ERR_SERVICE_UNAVAILABLE", "Service temporarily unavailable", http.StatusServiceUnavailable, nil)

	ErrJWTGeneration    = New("ERR_JWT_GENERATION", "Failed to generate JWT", http.StatusInternalServerError, nil)
	ErrInvalidJWTToken  = New("ERR_INVALID_JWT_TOKEN", "Invalid JWT token", http.StatusUnauthorized, nil)
	ErrJWTInvalidClaims = New("ERR_JWT_INVALID_CLAIM", "Invalid JWT claim", http.StatusUnauthorized, nil)
	ErrJWTExpired       = New("ERR_JWT_EXPIRED", "JWT token expired", http.StatusUnauthorized, nil)
)
