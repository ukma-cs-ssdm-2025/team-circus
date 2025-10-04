package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ukma-cs-ssdm-2025/team-circus/internal/config"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 20 * time.Second
)

type App struct {
	cfg *config.Config
	db  *sql.DB
	l   *zap.Logger
}

func New(cfg *config.Config, l *zap.Logger) *App {
	return &App{
		cfg: cfg,
		l:   l,
	}
}

func (a *App) Run(ctx context.Context) error {
	var err error

	a.db, err = sql.Open(a.cfg.DB.Driver, a.cfg.DB.DSN())
	if err != nil {
		return err
	}
	a.l.Info("DB connected")

	// wait for shutdown signal
	<-ctx.Done()
	a.l.Info("Shutdown signal received")

	// graceful shutdown
	timeoutCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	a.l.Info("Shutting down...")
	err = a.shutdown(timeoutCtx)
	if err != nil {
		a.l.Error("Shutdown error", zap.Error(err))
	}
	return err
}

func (a *App) shutdown(timeoutCtx context.Context) error {
	var shutdownErr error

	// db
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			wrapped := fmt.Errorf("shutdown db: %w", err)
			log.Println(wrapped)
			if shutdownErr == nil {
				shutdownErr = wrapped
			}
		} else {
			a.l.Info("DB closed")
		}
	}
	return shutdownErr
}
