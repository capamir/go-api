package database

import (
	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/pkg/logger"
)

// AutoMigrate runs automatic migrations for all models
func AutoMigrate() error {
	logger.Info("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		// We'll add more models here later
		// &models.Product{},
		// &models.Order{},
		// &models.OrderItem{},
	)

	if err != nil {
		logger.Error("Migration failed: %v", err)
		return err
	}

	logger.Success("Database migrations completed successfully!")
	return nil
}
