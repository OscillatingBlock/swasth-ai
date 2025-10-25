package server

import (
	"context"
	"net/http"
	"swasthAI/internal/auth/models"
	"swasthAI/internal/auth/repository"
	"swasthAI/internal/auth/usecase"
	"swasthAI/internal/middleware"

	authHandler "swasthAI/internal/auth/delivery/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	//init repos
	authRepo := repository.NewUserRepository(s.db, *s.logger)
	otpRepo := repository.NewOTPRepository(s.db)

	//init usecases
	authUC := usecase.NewAuthUsecase(authRepo, otpRepo, *s.cfg, *s.logger)

	//init handlers
	authHandler := authHandler.NewHandler(authUC, s.logger, s.cfg)

	//create tables
	ctx := context.Background()
	if _, err := s.db.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx); err != nil {
		s.logger.Error(err)
	}
	if _, err := s.db.NewCreateTable().Model((*models.OTP)(nil)).IfNotExists().Exec(ctx); err != nil {
		s.logger.Error(err)
	}

	//init middleware
	mw := middleware.NewMiddlewareManager(authUC, *s.cfg, s.logger)
	e.Use(mw.LoggerMiddleware)
	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	authHandler.MapAuthRoutes(authGroup, *mw)

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	return nil
}
