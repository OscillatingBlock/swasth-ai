package http

import (
	"github.com/labstack/echo/v4"

	"swasthAI/internal/middleware"
)

func (h *Handler) MapAuthRoutes(auth *echo.Group, mw middleware.MiddlewareManager) {
	auth.POST("/send-otp", h.SendOTP)
	auth.POST("/verify-otp", h.VerifyOTP)
	auth.POST("/register", h.Register)
	auth.POST("/refresh", h.RefreshToken)

	profileGroup := auth.Group("/profile")
	profileGroup.Use(mw.AuthJWTMiddleware)
	profileGroup.GET("", h.GetProfile)
	profileGroup.PUT("", h.UpdateProfile)
}
