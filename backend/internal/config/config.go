package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTMaxAge      time.Duration
	AdminEmail     string
	AdminPassword  string
	AllowedOrigins []string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getenv("PORT", "8080"),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTMaxAge:      parseDuration(getenv("JWT_TTL", "24h")),
		AdminEmail:     os.Getenv("ADMIN_EMAIL"),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
		AllowedOrigins: splitList(os.Getenv("ALLOWED_ORIGINS")),
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.JWTSecret == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if cfg.AdminEmail == "" {
		missing = append(missing, "ADMIN_EMAIL")
	}
	if cfg.AdminPassword == "" {
		missing = append(missing, "ADMIN_PASSWORD")
	}
	if cfg.AdminPassword != "" && len(cfg.AdminPassword) < 8 {
		return nil, fmt.Errorf("ADMIN_PASSWORD must be at least 8 characters")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	return 24 * time.Hour
}

func splitList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func (c *Config) PortInt() int {
	n, _ := strconv.Atoi(c.Port)
	return n
}
