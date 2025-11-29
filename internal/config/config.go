package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Server   ServerConfig
}

// AppConfig holds application-specific settings
type AppConfig struct {
	Environment string
	Port        string
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// JWTConfig holds JWT authentication settings
type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Global config instance - RENAMED to avoid conflict!
var Cfg *Config

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file (ignore error in production where env vars are set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Validate required variables
	requiredVars := []string{"JWT_SECRET", "DB_PASSWORD"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", v)
		}
	}

	config := &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "ecommerce"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: time.Hour * time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)),
		},
		Server: ServerConfig{
			ReadTimeout:  time.Second * time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 15)),
			WriteTimeout: time.Second * time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 15)),
		},
	}

	// Set global config
	Cfg = config

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as integer or returns a default
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsDevelopment checks if app is in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction checks if app is in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}
