package repository

import (
	"github.com/capamir/go-api/internal/models"
	"gorm.io/gorm"
)

// ProductRepository handles product database operations
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create inserts a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID retrieves a product by ID with category and tags
func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Tags").First(&product, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetBySlug retrieves a product by slug with category and tags
func (r *ProductRepository) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").Preload("Tags").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// GetBySKU retrieves a product by SKU
func (r *ProductRepository) GetBySKU(sku string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

// ProductFilterOptions represents filter options for products
type ProductFilterOptions struct {
	CategoryID  *uint
	TagIDs      []uint
	MinPrice    *float64
	MaxPrice    *float64
	InStock     *bool
	IsActive    *bool
	IsFeatured  *bool
	SearchQuery string
}

// GetAll retrieves all products with pagination and filters
func (r *ProductRepository) GetAll(limit, offset int, filters *ProductFilterOptions) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Apply filters
	if filters != nil {
		if filters.CategoryID != nil {
			query = query.Where("category_id = ?", *filters.CategoryID)
		}

		if len(filters.TagIDs) > 0 {
			query = query.Joins("JOIN product_tags ON product_tags.product_id = products.id").
				Where("product_tags.tag_id IN ?", filters.TagIDs).
				Group("products.id")
		}

		if filters.MinPrice != nil {
			query = query.Where("price >= ?", *filters.MinPrice)
		}

		if filters.MaxPrice != nil {
			query = query.Where("price <= ?", *filters.MaxPrice)
		}

		if filters.InStock != nil && *filters.InStock {
			query = query.Where("quantity > 0")
		}

		if filters.IsActive != nil {
			query = query.Where("is_active = ?", *filters.IsActive)
		}

		if filters.IsFeatured != nil {
			query = query.Where("is_featured = ?", *filters.IsFeatured)
		}

		if filters.SearchQuery != "" {
			searchPattern := "%" + filters.SearchQuery + "%"
			query = query.Where("name LIKE ? OR description LIKE ? OR sku LIKE ?", 
				searchPattern, searchPattern, searchPattern)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results with preloading
	err := query.Preload("Category").Preload("Tags").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetFeaturedProducts retrieves featured products
func (r *ProductRepository) GetFeaturedProducts(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Preload("Tags").
		Where("is_featured = ? AND is_active = ?", true, true).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error
	return products, err
}

// GetRelatedProducts retrieves products from the same category
func (r *ProductRepository) GetRelatedProducts(productID, categoryID uint, limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Category").Preload("Tags").
		Where("category_id = ? AND id != ? AND is_active = ?", categoryID, productID, true).
		Limit(limit).
		Order("RAND()").
		Find(&products).Error
	return products, err
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// UpdateTags updates product tags (many-to-many relationship)
func (r *ProductRepository) UpdateTags(product *models.Product, tags []models.Tag) error {
	return r.db.Model(product).Association("Tags").Replace(tags)
}

// Delete soft deletes a product
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// SlugExists checks if a slug already exists
func (r *ProductRepository) SlugExists(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("slug = ?", slug).Count(&count).Error
	return count > 0, err
}

// SKUExists checks if a SKU already exists
func (r *ProductRepository) SKUExists(sku string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Product{}).Where("sku = ?", sku).Count(&count).Error
	return count > 0, err
}

// UpdateStock updates product quantity
func (r *ProductRepository) UpdateStock(productID uint, quantity int) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).Update("quantity", quantity).Error
}

// DecrementStock decrements product quantity (for orders)
func (r *ProductRepository) DecrementStock(productID uint, amount int) error {
	return r.db.Model(&models.Product{}).
		Where("id = ? AND quantity >= ?", productID, amount).
		UpdateColumn("quantity", gorm.Expr("quantity - ?", amount)).Error
}
