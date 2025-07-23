package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort      string `envconfig:"HTTP_PORT" default:"8080"`
	DBHost        string `envconfig:"DB_HOST" default:"localhost"`
	DBPort        string `envconfig:"DB_PORT" default:"5432"`
	DBUser        string `envconfig:"DB_USER" default:"postgres"`
	DBPassword    string `envconfig:"DB_PASSWORD" default:"BIGsecret"`
	DBName        string `envconfig:"DB_NAME" default:"subscriptions_db"`
	DBSSLMode     string `envconfig:"DB_SSLMODE" default:"disable"`
	LogLevel      string `envconfig:"LOG_LEVEL" default:"debug"`
	MigrationPath string `envconfig:"MIGRATION_PATH" default:"file://migrations"`
}

func MustLoad() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to load config from env err: %v", err)
	}
	return cfg
}
