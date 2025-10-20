package testapp

import (
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/app"
	"github.com/ukma-cs-ssdm-2025/team-circus/internal/config"
	"go.uber.org/zap"
)

const Addr = "http://localhost:8081"

func NewApp() *app.App {
	cfg := config.Config{
		DB: config.DBConfig{
			Driver: "postgres",
			Host:   "localhost",
			Port:   "5433",
			User:   "postgres",
			Pass:   "postgres",
			Name:   "mcd",
		},
		Srv: config.SrvConfig{
			Port: "8081",
		},
		CORS: config.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			MaxAge:           43200,
		},
	}
	return app.New(&cfg, zap.NewNop())
}
