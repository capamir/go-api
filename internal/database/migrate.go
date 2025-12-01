package database

import (
	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/pkg/logger"
)

// AutoMigrate runs database migrations
func AutoMigrate() error {
	logger.Info("Running database migrations...")

	// Add all models here
	err := DB.AutoMigrate(
		&models.User{},
		&models.Category{}, // 🆕 Add this
	)

	if err != nil {
		return err
	}

	logger.Success("Database migrations completed successfully")
	return nil
}
