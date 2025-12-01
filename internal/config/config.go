package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
	Email    EmailConfig // 🆕 Add this
}

type AppConfig struct {
	Name string
	Env  string
	Port string
	URL  string // 🆕 Add this
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// 🆕 Add this struct
type EmailConfig struct {
	From     string
	Password string
	Host     string
	Port     int
	UseTLS   bool
}

func Load() (*Config, error) {
	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "E-Commerce API"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			URL:  getEnv("APP_URL", "http://localhost:8080"), // 🆕 Add this
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "ecommerce"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "default-secret-change-this"),
			Expiration: parseDuration(getEnv("JWT_EXPIRATION", "24h")),
		},
		Server: ServerConfig{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		// 🆕 Add this
		Email: EmailConfig{
			From:     getEnv("EMAIL_FROM", ""),
			Password: getEnv("EMAIL_PASSWORD", ""),
			Host:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
			Port:     getEnvInt("EMAIL_PORT", 587),
			UseTLS:   getEnvBool("EMAIL_USE_TLS", true),
		},
	}

	return cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	return d
}
