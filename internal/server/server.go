package server

import (
	"swasthAI/config"
	"swasthAI/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type Server struct {
	e      *echo.Echo
	logger *logger.Logger
	db     *bun.DB
	cfg    *config.Config
}

func NewServer(logger *logger.Logger, db *bun.DB, cfg *config.Config) *Server {
	return &Server{e: echo.New(), logger: logger, db: db, cfg: cfg}
}

func (s *Server) Run() error {
	s.MapHandlers(s.e)
	return s.e.Start(s.cfg.Server.Port)
}
