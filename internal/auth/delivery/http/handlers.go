package http

import (
	"swasthAI/config"
	"swasthAI/internal/auth/models"
	"swasthAI/internal/auth/usecase"
	appErrors "swasthAI/pkg/errors"
	"swasthAI/pkg/http_errors"
	"swasthAI/pkg/logger"
	"swasthAI/pkg/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	uc     *usecase.AuthUsecase
	logger *logger.Logger
	Cfg    *config.Config
}

func NewHandler(uc *usecase.AuthUsecase, logger *logger.Logger, cfg *config.Config) *Handler {
	return &Handler{uc: uc, logger: logger, Cfg: cfg}
}

func (h *Handler) Register(c echo.Context) error {
	var user models.User
	if err := utils.ReadRequest(c, &user); err != nil {
		h.logger.Error("Error while reading request", "error", err)
		return http_errors.Send(c, appErrors.ErrInvalidInput)
	}
	userWithToken, err := h.uc.Register(c.Request().Context(), &user)
	if err != nil {
		return err
	}
	return c.JSON(200, userWithToken)
}

func (h *Handler) Login(uc *usecase.AuthUsecase, c echo.Context) error {
	var user models.User
	if err := utils.ReadRequest(c, &user); err != nil {
		h.logger.Error(err)
		return http_errors.Send(c, appErrors.ErrInvalidInput)
	}
	userWithToken, err := uc.Login(c.Request().Context(), &user)
	if err != nil {
		return err
	}
	return c.JSON(200, userWithToken)
}

// NOTE: Using mock otp virification for now
func (h *Handler) VerifyOTP(c echo.Context) error {
	var user models.User
	if err := utils.ReadRequest(c, &user); err != nil {
		h.logger.Error(err)
		return http_errors.Send(c, appErrors.ErrInvalidInput)
	}
	return c.JSON(200, nil)
}
