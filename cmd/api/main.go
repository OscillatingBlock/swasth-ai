package main

import (
	"log/slog"
	"swasthAI/config"
	"swasthAI/internal/server"
	"swasthAI/pkg/logger"

	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	RawConfig, err := config.LoadConfig("config.yaml")
	if err != nil {
		slog.Error("error while loading config", "err", err)
	}
	cfg, err := config.ParseConfig(RawConfig)
	if err != nil {
		slog.Error("error while parsing config", "err", err)
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.Bun.DSN)))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.PingContext(context.Background()); err != nil {
		slog.Error("❌ failed to connect to Postgres: %v", "err", err)
	}

	slog.Info("✅ Connected to Postgres successfully")

	defer db.Close()

	// Test connection
	if err := db.PingContext(context.Background()); err != nil {
		slog.Error("error while pinging db", "err", err)
	}

	slog.Info("database connected successfully!")

	logger, err := logger.NewLogger(cfg)
	if err != nil {
		slog.Error("error while getting logger", "err", err)
	}

	server := server.NewServer(logger, db, cfg)
	server.Run()
}
