// @title Team Circus API
// @version 1.0
// @description API for Team Circus project
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"context"
	"os/signal"
	"syscall"

	_ "github.com/swaggo/gin-swagger"
	_ "github.com/ukma-cs-ssdm-2025/team-circus/docs"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/app"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/config"
	"github.com/ukma-cs-ssdm-2025/team-circus/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	logger := logging.NewLogger()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	cfg, err := config.Load()
	if err != nil {
		logger.Panic("Failed to load config", zap.Error(err))
	}
	logger.Info("Creating app...")
	app := app.New(cfg, logger)
	logger.Info("App created")

	logger.Info("Running app...")
	if err := app.Run(ctx); err != nil {
		logger.Panic("Failed to run app", zap.Error(err))
	}
	logger.Info("App stopped successfully")
}
