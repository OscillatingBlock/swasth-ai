package middleware

import (
	"swasthAI/config"
	"swasthAI/internal/auth/usecase"
	"swasthAI/pkg/logger"
)

type MiddlewareManager struct {
	AuthUC usecase.AuthUsecase
	Cfg    config.Config
	Logger *logger.Logger
}

func NewMiddlewareManager(uc *usecase.AuthUsecase, cfg config.Config, logger *logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{AuthUC: *uc, Cfg: cfg, Logger: logger}
}
