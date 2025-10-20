package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/config"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 20 * time.Second
	readTimeout     = 15 * time.Second
)

type App struct {
	cfg *config.Config
	DB  *sql.DB
	API *http.Server
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

	// db
	a.DB, err = sql.Open(a.cfg.DB.Driver, a.cfg.DB.DSN())
	if err != nil {
		return err
	}
	a.l.Info("DB connected")

	// api
	router := a.setupRouter()
	a.API = &http.Server{
		Addr:        ":" + a.cfg.Srv.Port,
		Handler:     router,
		ReadTimeout: readTimeout,
	}
	go func() {
		if err := a.API.ListenAndServe(); err != nil {
			log.Printf("api server: %v", err)
		}
	}()
	log.Printf("APIServer started on port %s", a.cfg.Srv.Port)

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

	// api
	if a.API != nil {
		if err := a.API.Shutdown(timeoutCtx); err != nil {
			wrapped := fmt.Errorf("shutdown api server: %w", err)
			log.Println(wrapped)
			shutdownErr = wrapped
		} else {
			log.Println("APIServer Shutdown successfully")
		}
	}

	// db
	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
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
