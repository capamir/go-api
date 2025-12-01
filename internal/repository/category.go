package repository

import (
	"github.com/capamir/go-api/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository handles category database operations
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create inserts a new category
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found, not an error
		}
		return nil, err
	}
	return &category, nil
}

// GetBySlug retrieves a category by slug
func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetAll retrieves all categories with pagination
func (r *CategoryRepository) GetAll(limit, offset int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total
	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Limit(limit).Offset(offset).Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// GetRootCategories retrieves categories without parent (top-level)
func (r *CategoryRepository) GetRootCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("parent_id IS NULL").Order("name ASC").Find(&categories).Error
	return categories, err
}

// GetWithChildren retrieves a category with all its children
func (r *CategoryRepository) GetWithChildren(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Children").First(&category, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// GetCategoryTree retrieves all categories with nested children
func (r *CategoryRepository) GetCategoryTree() ([]models.Category, error) {
	// Get all root categories
	var rootCategories []models.Category
	err := r.db.Preload("Children").Where("parent_id IS NULL").Order("name ASC").Find(&rootCategories).Error
	return rootCategories, err
}

// Update updates a category
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete soft deletes a category
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// SlugExists checks if a slug already exists
func (r *CategoryRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

// GetActiveCategories retrieves all active categories
func (r *CategoryRepository) GetActiveCategories(limit, offset int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total active
	if err := r.db.Model(&models.Category{}).Where("is_active = ?", true).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Where("is_active = ?", true).Limit(limit).Offset(offset).Order("name ASC").Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}
