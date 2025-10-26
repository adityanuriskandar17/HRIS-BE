package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBURL         string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
	Telemetry     TelemetryConfig
}

type TelemetryConfig struct {
	ServiceName    string
	JaegerEndpoint string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Port:          get("PORT", "8080"),
		DBURL:         get("DATABASE_URL", "postgres://hris:hris@127.0.0.1:5432/hris?sslmode=disable"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AdminEmail:    get("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: get("ADMIN_PASSWORD", "admin12345"),
		Telemetry: TelemetryConfig{
			ServiceName:    get("OTEL_SERVICE_NAME", "hris-api"),
			JaegerEndpoint: get("OTEL_EXPORTER_JAEGER_ENDPOINT", ""),
		},
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET required")
	}
	return cfg, nil
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
