package db

import (
	"database/sql"
	"fmt"
	"time"
	"unicode"

	"github.com/capamir/go-api/utils"
	"github.com/go-sql-driver/mysql"
)

// holds database connection pool settings
type ConnectionConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// returns recommended connection pool settings
func DefaultConnectionConfig() ConnectionConfig {
	return ConnectionConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// creates a new MySQL database connection with auto-creation support
func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	return NewMySQLStorageWithConfig(cfg, DefaultConnectionConfig())
}

// creates a new MySQL connection with custom pool config
func NewMySQLStorageWithConfig(cfg mysql.Config, connCfg ConnectionConfig) (*sql.DB, error) {
	utils.S.DBStatus("connecting", "Connecting to MySQL...")

	// Validate database name before proceeding
	if !isValidDatabaseName(cfg.DBName) {
		return nil, fmt.Errorf("invalid database name: %s (only alphanumeric and underscores allowed)", cfg.DBName)
	}

	// Phase 1: Connect without database to create it if needed
	cfgWithoutDB := cfg
	cfgWithoutDB.DBName = ""
	
	db, err := sql.Open("mysql", cfgWithoutDB.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database server: %w", err)
	}

	// Create database if it doesn't exist
	if err := createDatabaseIfNotExists(db, cfg.DBName); err != nil {
		db.Close()
		return nil, err
	}

	// Close initial connection
	db.Close()

	// Phase 2: Connect to the specific database
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database %s: %w", cfg.DBName, err)
	}

	// Verify connection to specific database
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database %s: %w", cfg.DBName, err)
	}

	// Configure connection pool
	configureConnectionPool(db, connCfg)

	utils.S.DBStatus("connected", fmt.Sprintf("Connected to database '%s'", cfg.DBName))
	return db, nil
}

// createDatabaseIfNotExists creates the database if it doesn't already exist
func createDatabaseIfNotExists(db *sql.DB, dbName string) error {
	// Double-check validation (defense in depth)
	if !isValidDatabaseName(dbName) {
		return fmt.Errorf("invalid database name: %s", dbName)
	}

	// Note: MySQL doesn't support parameterized database names
	// We use backticks and validation to prevent SQL injection
	createQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	
	if _, err := db.Exec(createQuery); err != nil {
		return fmt.Errorf("failed to create database %s: %w", dbName, err)
	}

	utils.S.Debugf("Database '%s' created or already exists", dbName)

	// Select the database
	useQuery := fmt.Sprintf("USE `%s`", dbName)
	if _, err := db.Exec(useQuery); err != nil {
		return fmt.Errorf("failed to select database %s: %w", dbName, err)
	}

	return nil
}

// isValidDatabaseName validates that the database name contains only safe characters
func isValidDatabaseName(name string) bool {
	// Check length (MySQL limit is 64 characters)
	if len(name) == 0 || len(name) > 64 {
		return false
	}

	// Only allow alphanumeric characters and underscores
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false
		}
	}

	return true
}

// configureConnectionPool sets up the connection pool parameters
func configureConnectionPool(db *sql.DB, cfg ConnectionConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	utils.S.Debugf(
		"Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v",
		cfg.MaxOpenConns,
		cfg.MaxIdleConns,
		cfg.ConnMaxLifetime,
	)
}

// TestConnection tests the database connection with retries
func TestConnection(db *sql.DB, retries int, retryDelay time.Duration) error {
	var err error
	
	for i := 0; i < retries; i++ {
		err = db.Ping()
		if err == nil {
			utils.S.Successf("Database connection verified (attempt %d/%d)", i+1, retries)
			return nil
		}
		
		utils.S.Warnf("Database ping attempt %d/%d failed: %v", i+1, retries, err)
		
		if i < retries-1 {
			time.Sleep(retryDelay)
		}
	}
	
	return fmt.Errorf("failed to connect after %d attempts: %w", retries, err)
}

// GetDatabaseStats returns connection pool statistics
func GetDatabaseStats(db *sql.DB) sql.DBStats {
	return db.Stats()
}

// LogDatabaseStats logs current database connection pool statistics
func LogDatabaseStats(db *sql.DB) {
	stats := db.Stats()
	utils.S.Infof(
		"DB Stats - OpenConns: %d, InUse: %d, Idle: %d, WaitCount: %d, WaitDuration: %v",
		stats.OpenConnections,
		stats.InUse,
		stats.Idle,
		stats.WaitCount,
		stats.WaitDuration,
	)
}
