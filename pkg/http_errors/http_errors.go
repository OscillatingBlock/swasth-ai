package http_errors

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"swasthAI/internal/auth/models"
	"swasthAI/pkg/domain_errors"
	appErrors "swasthAI/pkg/errors"
)

// ErrorResponse standardizes HTTP error responses
type ErrorResponse struct {
	Error      string                 `json:"error"`
	Code       string                 `json:"code"`
	Details    map[string]interface{} `json:"details,omitempty"`
	RetryAfter int                    `json:"retry_after,omitempty"`
}

// Send sends standardized error response
func Send(c echo.Context, appErr *appErrors.AppError) error {
	details := map[string]interface{}{}

	// Add specific details based on error code
	switch appErr.Code {
	case "USER_INVALID_LANGUAGE":
		details["supported"] = []string{"hi", "en", "ta", "te", "bn", "mr"}
	case "VIDEO_INVALID_CATEGORY":
		details["valid_categories"] = []string{"snake_bite", "cpr", "burns", "bleeding"}
	case "VOICE_INVALID_FORMAT":
		details["supported"] = []string{"WAV", "MP3"}
	case "VISION_INVALID_IMAGE":
		details["supported"] = []string{"JPEG", "PNG"}
	case "VOICE_AUDIO_TOO_LARGE":
		details["max_size"] = "10MB"
	case "VISION_IMAGE_TOO_LARGE":
		details["max_size"] = "5MB"
	case "VISION_PDF_TOO_LARGE":
		details["max_size"] = "10MB"
	case "AUTH_INVALID_OTP":
		details["retry_attempts"] = 3
	case "AUTH_RESEND_COOLDOWN":
		details["retry_after"] = 60
	case "ERR_RATE_LIMITED":
		details["retry_after"] = 60
	}

	retryAfter := 0
	if appErr.Status == http.StatusTooManyRequests {
		retryAfter = 60
	}

	resp := ErrorResponse{
		Error:      appErr.Message,
		Code:       appErr.Code,
		Details:    details,
		RetryAfter: retryAfter,
	}

	return c.JSON(appErr.Status, resp)
}

// Common HTTP Error Handlers
func Handle(c echo.Context, err error) error {
	// Extract AppError
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		return Send(c, appErr)
	}

	// Handle Bun ORM errors
	if err == sql.ErrNoRows {
		return Send(c, appErrors.ErrDatabase)
	}

	// Default internal error
	return Send(c, appErrors.ErrInternal)
}

// ValidationErrorResponse for Echo validator
func ValidationErrorResponse(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
		"error":   "Validation failed",
		"code":    "ERR_VALIDATION",
		"details": err.Error(),
	})
}

// User-Specific Validation Helper
func ValidateUserRequest(c echo.Context, user *models.User) error {
	if err := c.Validate(user); err != nil {
		return ValidationErrorResponse(err)
	}

	// Custom domain validation
	if err := domain_errors.ValidateUserPhone(user.Phone); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
			"code":  err.(*appErrors.AppError).Code,
		})
	}

	if err := domain_errors.ValidateUserLanguage(user.Language); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, map[string]interface{}{
			"error":     err.Error(),
			"code":      err.(*appErrors.AppError).Code,
			"supported": []string{"hi", "en", "ta", "te", "bn", "mr"},
		})
	}

	return nil
}
