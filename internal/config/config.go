package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DBURL         string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
	Auth          AuthConfig
	CORS          CORSConfig
	Telemetry     TelemetryConfig
}

type TelemetryConfig struct {
	ServiceName    string
	JaegerEndpoint string
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	AllowedOrigins  []string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	accessTTL, err := parseBoundedDuration(get("ACCESS_TOKEN_TTL", "15m"), 15*time.Minute, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := parseBoundedDuration(get("REFRESH_TOKEN_TTL", "336h"), 7*24*time.Hour, 30*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("REFRESH_TOKEN_TTL: %w", err)
	}

	// Parse CORS configuration
	corsAllowedOrigins := parseCSV(get("CORS_ALLOWED_ORIGINS", ""))
	corsAllowedMethods := parseCSV(get("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"))
	corsAllowedHeaders := parseCSV(get("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"))
	corsExposedHeaders := parseCSV(get("CORS_EXPOSED_HEADERS", "Link"))
	corsAllowCredentials := get("CORS_ALLOW_CREDENTIALS", "true") == "true"
	corsMaxAge := parseInt(get("CORS_MAX_AGE", "300"), 300)

	// Default origins for development if not specified
	if len(corsAllowedOrigins) == 0 {
		corsAllowedOrigins = []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:5173", // Vite default
			"http://localhost:8080",
			"http://localhost:8081",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:5173",
		}
	}

	cfg := &Config{
		Port:          get("PORT", "8080"),
		DBURL:         get("DATABASE_URL", "postgres://hris:hris@127.0.0.1:5432/hris?sslmode=disable"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AdminEmail:    get("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: get("ADMIN_PASSWORD", "admin12345"),
		Auth: AuthConfig{
			AccessTokenTTL:  accessTTL,
			RefreshTokenTTL: refreshTTL,
			AllowedOrigins:  corsAllowedOrigins,
		},
		CORS: CORSConfig{
			AllowedOrigins:   corsAllowedOrigins,
			AllowedMethods:   corsAllowedMethods,
			AllowedHeaders:   corsAllowedHeaders,
			ExposedHeaders:   corsExposedHeaders,
			AllowCredentials: corsAllowCredentials,
			MaxAge:           corsMaxAge,
		},
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

func parseCSV(input string) []string {
	if input == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseBoundedDuration(raw string, min, max time.Duration) (time.Duration, error) {
	v, err := time.ParseDuration(raw)
	if err != nil {
		return 0, err
	}
	if v < min || v > max {
		return 0, fmt.Errorf("duration %s outside allowed range [%s, %s]", v, min, max)
	}
	return v, nil
}

func parseInt(raw string, defaultValue int) int {
	var result int
	if _, err := fmt.Sscanf(raw, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}
