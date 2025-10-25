package http

import (
	"swasthAI/config"
	"swasthAI/internal/auth"
	"swasthAI/internal/auth/models"
	"swasthAI/internal/auth/usecase"
	appErrors "swasthAI/pkg/errors"
	"swasthAI/pkg/http_errors"
	"swasthAI/pkg/logger"
	"swasthAI/pkg/utils"

	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type Handler struct {
	uc     auth.AuthUsecase
	logger *logger.Logger
	Cfg    *config.Config
}

func NewHandler(uc *usecase.AuthUsecase, logger *logger.Logger, cfg *config.Config) *Handler {
	return &Handler{uc: uc, logger: logger, Cfg: cfg}
}

func (h *Handler) SendOTP(c echo.Context) error {
	var user models.SendOTPInput
	if err := utils.ReadRequest(c, &user); err != nil {
		h.logger.Error(err)
		return http_errors.Send(c, appErrors.ErrInvalidInput)
	}
	err := h.uc.SendOTP(c.Request().Context(), user.Phone)
	if err != nil {
		h.logger.Error("failed to send otp", "error", err)

		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInternal)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "OTP sent to " + user.Phone,
		"retry_after": "60",
	})
}

func (h *Handler) Register(c echo.Context) error {
	var input models.RegisterUserInput
	if err := utils.ReadRequest(c, &input); err != nil {
		h.logger.Error("failed to read input", "error", err)

		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInvalidInput)
	}

	userWithToken, err := h.uc.RegisterUser(c.Request().Context(), &input)
	if err != nil {
		h.logger.Error("failed to register user", "error", err)

		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrDatabase)
	}

	return c.JSON(http.StatusOK, userWithToken)
}

// NOTE: Using mock otp verification for now
func (h *Handler) VerifyOTP(c echo.Context) error {
	var input models.VerifyOTPInput
	if err := utils.ReadRequest(c, &input); err != nil {
		h.logger.Error("failed to read input", "error", err)

		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInvalidInput)
	}

	userWithToken, registered, err := h.uc.VerifyOTP(c.Request().Context(), input.Phone, input.OTP)
	if err != nil {
		h.logger.Error("failed to verify otp", "error", err)

		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrDatabase)
	}

	// Success: User exists → Login
	if registered {
		return c.JSON(http.StatusOK, userWithToken)
	}

	// Success: User does not exist → Redirect to register
	return c.JSON(http.StatusOK, map[string]any{
		"registered": false,
		"message":    "Complete registration",
	})
}

// func (h *Handler) RefreshToken(c echo.Context) error {
// 	token, err := h.uc.RefreshToken(c.Request().Context(), c.Request().Header.Get("Authorization"))
// 	if err != nil {
// 		h.logger.Error("failed to refresh token", "error", err)

// 		if appErr, ok := err.(*appErrors.AppError); ok {
// 			return http_errors.Send(c, appErr)
// 		}
// 		return http_errors.Send(c, appErrors.ErrInternal)
// 	}

// 	return c.JSON(http.StatusOK, map[string]any{
// 		"token":         token.Token,
// 		"refresh_token": token.RefreshToken,
// 		"expires_in":    86400,
// 	})
// }

func (h *Handler) RefreshToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		h.logger.Error("missing Authorization header")
		return http_errors.Send(c, appErrors.ErrUnauthorized)
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		h.logger.Error("invalid Authorization header format")
		return http_errors.Send(c, appErrors.ErrUnauthorized)
	}

	refreshToken := parts[1]

	token, err := h.uc.RefreshToken(c.Request().Context(), refreshToken)
	if err != nil {
		h.logger.Error("failed to refresh token", "error", err)
		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInternal)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token":         token.Token,
		"refresh_token": token.RefreshToken,
		"expires_in":    86400,
	})
}

func (h *Handler) GetProfile(c echo.Context) error {
	user, err := h.uc.GetUserByID(c.Request().Context(), c.Request().Header.Get("Authorization"))
	if err != nil {
		h.logger.Error("failed to get user by id", "error", err)

		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInternal)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	var input models.UpdateProfileInput
	if err := utils.ReadRequest(c, &input); err != nil {
		h.logger.Error("failed to read input", "error", err)
		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInternal)
	}

	_, err := h.uc.UpdateProfile(c.Request().Context(), &input)
	if err != nil {
		h.logger.Error("failed to update profile", "error", err)
		if appErr, ok := err.(*appErrors.AppError); ok {
			return http_errors.Send(c, appErr)
		}
		return http_errors.Send(c, appErrors.ErrInternal)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Profile updated successfully",
	})
}
