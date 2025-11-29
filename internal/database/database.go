package database

import (
	"fmt"
	"time"

	"github.com/capamir/go-api/internal/config"
	"github.com/capamir/go-api/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Connect establishes a connection to the MySQL database
func Connect(cfg *config.Config) error {
	logger.Log.DB("connecting", "Connecting to MySQL database...")

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// Configure GORM
	var logLevel gormLogger.LogLevel
	if cfg.IsDevelopment() {
		logLevel = gormLogger.Info
	} else {
		logLevel = gormLogger.Error
	}

	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	}

	// Open connection
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		logger.Log.DB("error", fmt.Sprintf("Failed to connect: %v", err))
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying *sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		logger.Log.DB("error", fmt.Sprintf("Ping failed: %v", err))
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set global DB
	DB = db

	logger.Log.DB("connected", fmt.Sprintf("Successfully connected to database '%s'", cfg.Database.Name))
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
