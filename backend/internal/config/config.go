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

type SrvConfig struct {
	Port string `envconfig:"API_PORT" required:"true"`
}

type CORSConfig struct {
	AllowOrigins     []string `envconfig:"CORS_ALLOW_ORIGINS" required:"true"`
	AllowMethods     []string `envconfig:"CORS_ALLOW_METHODS" required:"true"`
	AllowHeaders     []string `envconfig:"CORS_ALLOW_HEADERS" required:"true"`
	ExposeHeaders    []string `envconfig:"CORS_EXPOSE_HEADERS" required:"true"`
	AllowCredentials bool     `envconfig:"CORS_ALLOW_CREDENTIALS" required:"true"`
	MaxAge           int      `envconfig:"CORS_MAX_AGE" required:"true"`
}

type ShareConfig struct {
	HMACSecret            string `envconfig:"SHARE_HMAC_SECRET" required:"true"`
	BaseURL               string `envconfig:"APP_BASE_URL" required:"true"`
	DefaultExpirationDays int    `envconfig:"SHARE_DEFAULT_EXPIRATION_DAYS" default:"7"`
	MaxExpirationDays     int    `envconfig:"SHARE_MAX_EXPIRATION_DAYS" default:"90"`
}

type Config struct {
	DB              DBConfig
	Srv             SrvConfig
	CORS            CORSConfig
	HashingCost     int    `envconfig:"HASHING_COST" required:"true"`
	SecretToken     string `envconfig:"SECRET_TOKEN" required:"true"`
	AccessDuration  int    `envconfig:"ACCESS_DURATION" required:"true"`
	RefreshDuration int    `envconfig:"REFRESH_DURATION" required:"true"`
	Share           ShareConfig
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
