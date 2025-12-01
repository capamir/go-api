package repository

import (
	"github.com/capamir/go-api/internal/models"
	"gorm.io/gorm"
)

// TagRepository handles tag database operations
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository creates a new tag repository
func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

// Create inserts a new tag
func (r *TagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// GetByID retrieves a tag by ID
func (r *TagRepository) GetByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found, not an error
		}
		return nil, err
	}
	return &tag, nil
}

// GetBySlug retrieves a tag by slug
func (r *TagRepository) GetBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// GetAll retrieves all tags with pagination
func (r *TagRepository) GetAll(limit, offset int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	// Count total
	if err := r.db.Model(&models.Tag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Limit(limit).Offset(offset).Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// GetAllTags retrieves all tags without pagination (for dropdowns, etc.)
func (r *TagRepository) GetAllTags() ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.Order("name ASC").Find(&tags).Error
	return tags, err
}

// Update updates a tag
func (r *TagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

// Delete soft deletes a tag
func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

// SlugExists checks if a slug already exists
func (r *TagRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Tag{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

// Search searches tags by name
func (r *TagRepository) Search(query string, limit, offset int) ([]models.Tag, int64, error) {
	var tags []models.Tag
	var total int64

	searchQuery := "%" + query + "%"

	// Count total matches
	if err := r.db.Model(&models.Tag{}).Where("name LIKE ?", searchQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Where("name LIKE ?", searchQuery).Limit(limit).Offset(offset).Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}
