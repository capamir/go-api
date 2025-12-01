package repository

import (
	"errors"

	"github.com/capamir/go-api/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found (not an error)
		}
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// GetAll retrieves all users
func (r *UserRepository) GetAll(limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ========================================
// 🆕 EMAIL VERIFICATION METHODS
// ========================================

// GetByVerificationToken retrieves a user by their verification token
func (r *UserRepository) GetByVerificationToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("verification_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Token not found (not an error)
		}
		return nil, err
	}
	return &user, nil
}

// GetUnverifiedUsers retrieves users who haven't verified their email
// Useful for: sending reminder emails, cleanup old unverified accounts
func (r *UserRepository) GetUnverifiedUsers(limit int) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("email_verified = ?", false).
		Limit(limit).
		Order("created_at ASC"). // Oldest first
		Find(&users).Error
	return users, err
}
