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
		&models.Category{}, 
		&models.Tag{},
		&models.Product{},
		&models.Cart{},      
		&models.CartItem{},  
		&models.Order{},     
		&models.OrderItem{}, 
	)

	if err != nil {
		return err
	}

	logger.Success("Database migrations completed successfully")
	return nil
}
