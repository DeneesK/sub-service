package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddr    string        `envconfig:"SERVER_ADDR" default:"localhost:8080"`
	TimeOut       time.Duration `envconfig:"TIMEOUT" default:"30s"`
	DBHost        string        `envconfig:"DB_HOST" default:"localhost"`
	DBPort        string        `envconfig:"DB_PORT" default:"5432"`
	DBUser        string        `envconfig:"DB_USER" default:"postgres"`
	DBPassword    string        `envconfig:"DB_PASSWORD" default:"BIGsecret"`
	DBName        string        `envconfig:"DB_NAME" default:"subscriptions_db"`
	DBSSLMode     string        `envconfig:"DB_SSLMODE" default:"disable"`
	LogLevel      string        `envconfig:"LOG_LEVEL" default:"debug"`
	MigrationPath string        `envconfig:"MIGRATION_PATH" default:"file://migrations"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env found, using environment variables")
	}
}

func MustLoad() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to load config from env err: %v", err)
	}
	return cfg
}
