package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Driver string `envconfig:"DB_DRIVER" required:"true"`
	Host   string `envconfig:"DB_HOST" required:"true"`
	Port   string `envconfig:"DB_PORT" required:"true"`
	User   string `envconfig:"DB_USER" required:"true"`
	Pass   string `envconfig:"DB_PASSWORD" required:"true"`
	Name   string `envconfig:"DB_NAME" required:"true"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Pass, c.Name,
	)
}

type Config struct {
	DB DBConfig
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
