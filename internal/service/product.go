package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"gorm.io/datatypes"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo  *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
	tagRepo      *repository.TagRepository
}

// NewProductService creates a new product service
func NewProductService(
	productRepo *repository.ProductRepository,
	categoryRepo *repository.CategoryRepository,
	tagRepo *repository.TagRepository,
) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

// ========================================
// DTOs (Data Transfer Objects)
// ========================================

// CreateProductRequest represents product creation data
type CreateProductRequest struct {
	Name         string   `json:"name" binding:"required,min=2,max=255"`
	Description  string   `json:"description" binding:"omitempty,max=5000"`
	Price        float64  `json:"price" binding:"required,gt=0"`
	ComparePrice float64  `json:"compare_price" binding:"omitempty,gtefield=Price"`
	SKU          string   `json:"sku" binding:"omitempty,max=100"`
	Barcode      string   `json:"barcode" binding:"omitempty,max=100"`
	Quantity     int      `json:"quantity" binding:"omitempty,gte=0"`
	Images       []string `json:"images" binding:"omitempty"`
	CategoryID   uint     `json:"category_id" binding:"required"`
	TagIDs       []uint   `json:"tag_ids" binding:"omitempty"`
	IsActive     *bool    `json:"is_active" binding:"omitempty"`
	IsFeatured   *bool    `json:"is_featured" binding:"omitempty"`
}

// UpdateProductRequest represents product update data
type UpdateProductRequest struct {
	Name         string   `json:"name" binding:"omitempty,min=2,max=255"`
	Description  string   `json:"description" binding:"omitempty,max=5000"`
	Price        float64  `json:"price" binding:"omitempty,gt=0"`
	ComparePrice float64  `json:"compare_price" binding:"omitempty"`
	SKU          string   `json:"sku" binding:"omitempty,max=100"`
	Barcode      string   `json:"barcode" binding:"omitempty,max=100"`
	Quantity     int      `json:"quantity" binding:"omitempty,gte=0"`
	Images       []string `json:"images" binding:"omitempty"`
	CategoryID   uint     `json:"category_id" binding:"omitempty"`
	TagIDs       []uint   `json:"tag_ids" binding:"omitempty"`
	IsActive     *bool    `json:"is_active" binding:"omitempty"`
	IsFeatured   *bool    `json:"is_featured" binding:"omitempty"`
}

// ProductFilterRequest represents product filter parameters
type ProductFilterRequest struct {
	CategoryID  *uint    `form:"category_id"`
	TagIDs      []uint   `form:"tag_ids"`
	MinPrice    *float64 `form:"min_price"`
	MaxPrice    *float64 `form:"max_price"`
	InStock     *bool    `form:"in_stock"`
	IsActive    *bool    `form:"is_active"`
	IsFeatured  *bool    `form:"is_featured"`
	SearchQuery string   `form:"q"`
	Page        int      `form:"page"`
	Limit       int      `form:"limit"`
}

// ProductListResponse represents paginated product list
type ProductListResponse struct {
	Products   []models.ProductResponse `json:"products"`
	Total      int64                    `json:"total"`
	Page       int                      `json:"page"`
	Limit      int                      `json:"limit"`
	TotalPages int                      `json:"total_pages"`
}

// ========================================
// CREATE PRODUCT
// ========================================

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(req CreateProductRequest) (*models.ProductResponse, error) {
	// Validate category exists
	category, err := s.categoryRepo.GetByID(req.CategoryID)
	if err != nil {
		logger.Error("Failed to check category: %v", err)
		return nil, errors.New("failed to create product")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Generate slug from name
	slug := utils.GenerateSlug(req.Name)

	// Check if slug already exists
	exists, err := s.productRepo.SlugExists(slug)
	if err != nil {
		logger.Error("Failed to check slug existence: %v", err)
		return nil, errors.New("failed to create product")
	}
	if exists {
		return nil, errors.New("product with similar name already exists")
	}

	// Check SKU uniqueness if provided
	if req.SKU != "" {
		exists, err := s.productRepo.SKUExists(req.SKU)
		if err != nil {
			logger.Error("Failed to check SKU existence: %v", err)
			return nil, errors.New("failed to create product")
		}
		if exists {
			return nil, errors.New("SKU already exists")
		}
	}

	// Validate and fetch tags
	var tags []models.Tag
	if len(req.TagIDs) > 0 {
		for _, tagID := range req.TagIDs {
			tag, err := s.tagRepo.GetByID(tagID)
			if err != nil {
				logger.Error("Failed to check tag %d: %v", tagID, err)
				return nil, errors.New("failed to create product")
			}
			if tag == nil {
				return nil, fmt.Errorf("tag with ID %d not found", tagID)
			}
			tags = append(tags, *tag)
		}
	}

	// Convert images to JSON
	var imagesJSON datatypes.JSON
	if len(req.Images) > 0 {
		imagesBytes, _ := json.Marshal(req.Images)
		imagesJSON = datatypes.JSON(imagesBytes)
	} else {
		imagesBytes, _ := json.Marshal([]string{})
		imagesJSON = datatypes.JSON(imagesBytes)
	}

	// Set default values
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	isFeatured := false
	if req.IsFeatured != nil {
		isFeatured = *req.IsFeatured
	}

	quantity := 0
	if req.Quantity > 0 {
		quantity = req.Quantity
	}

	// Create product
	product := &models.Product{
		Name:         req.Name,
		Slug:         slug,
		Description:  req.Description,
		Price:        req.Price,
		ComparePrice: req.ComparePrice,
		SKU:          req.SKU,
		Barcode:      req.Barcode,
		Quantity:     quantity,
		Images:       imagesJSON,
		CategoryID:   req.CategoryID,
		Tags:         tags,
		IsActive:     isActive,
		IsFeatured:   isFeatured,
	}

	if err := s.productRepo.Create(product); err != nil {
		logger.Error("Failed to create product: %v", err)
		return nil, errors.New("failed to create product")
	}

	logger.Success("Product created successfully: %s (id: %d, slug: %s)", product.Name, product.ID, product.Slug)

	// Load relationships for response
	product, _ = s.productRepo.GetByID(product.ID)
	response := product.ToResponse()
	return &response, nil
}

// ========================================
// GET PRODUCT
// ========================================

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(id uint) (*models.ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get product %d: %v", id, err)
		return nil, errors.New("failed to get product")
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	response := product.ToResponse()
	return &response, nil
}

// GetProductBySlug retrieves a product by slug
func (s *ProductService) GetProductBySlug(slug string) (*models.ProductResponse, error) {
	product, err := s.productRepo.GetBySlug(slug)
	if err != nil {
		logger.Error("Failed to get product by slug %s: %v", slug, err)
		return nil, errors.New("failed to get product")
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	response := product.ToResponse()
	return &response, nil
}

// ========================================
// LIST PRODUCTS
// ========================================

// GetAllProducts retrieves all products with filters and pagination
func (s *ProductService) GetAllProducts(filterReq ProductFilterRequest) (*ProductListResponse, error) {
	// Default pagination
	page := filterReq.Page
	if page < 1 {
		page = 1
	}

	limit := filterReq.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Build filters
	filters := &repository.ProductFilterOptions{
		CategoryID:  filterReq.CategoryID,
		TagIDs:      filterReq.TagIDs,
		MinPrice:    filterReq.MinPrice,
		MaxPrice:    filterReq.MaxPrice,
		InStock:     filterReq.InStock,
		IsActive:    filterReq.IsActive,
		IsFeatured:  filterReq.IsFeatured,
		SearchQuery: filterReq.SearchQuery,
	}

	products, total, err := s.productRepo.GetAll(limit, offset, filters)
	if err != nil {
		logger.Error("Failed to get products: %v", err)
		return nil, errors.New("failed to get products")
	}

	// Convert to response DTOs
	productResponses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = product.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &ProductListResponse{
		Products:   productResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetFeaturedProducts retrieves featured products
func (s *ProductService) GetFeaturedProducts(limit int) ([]models.ProductResponse, error) {
	if limit < 1 || limit > 50 {
		limit = 10
	}

	products, err := s.productRepo.GetFeaturedProducts(limit)
	if err != nil {
		logger.Error("Failed to get featured products: %v", err)
		return nil, errors.New("failed to get featured products")
	}

	// Convert to response DTOs
	productResponses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = product.ToResponse()
	}

	return productResponses, nil
}

// GetRelatedProducts retrieves related products (same category)
func (s *ProductService) GetRelatedProducts(productID uint, limit int) ([]models.ProductResponse, error) {
	// Get the product first to find its category
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		logger.Error("Failed to get product %d: %v", productID, err)
		return nil, errors.New("failed to get related products")
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if limit < 1 || limit > 20 {
		limit = 4
	}

	products, err := s.productRepo.GetRelatedProducts(productID, product.CategoryID, limit)
	if err != nil {
		logger.Error("Failed to get related products: %v", err)
		return nil, errors.New("failed to get related products")
	}

	// Convert to response DTOs
	productResponses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = product.ToResponse()
	}

	return productResponses, nil
}

// ========================================
// UPDATE PRODUCT
// ========================================

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(id uint, req UpdateProductRequest) (*models.ProductResponse, error) {
	// Get existing product
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get product %d: %v", id, err)
		return nil, errors.New("failed to update product")
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Update name and slug if provided
	if req.Name != "" && req.Name != product.Name {
		newSlug := utils.GenerateSlug(req.Name)

		// Check if new slug already exists (excluding current product)
		existingProduct, err := s.productRepo.GetBySlug(newSlug)
		if err != nil {
			logger.Error("Failed to check slug existence: %v", err)
			return nil, errors.New("failed to update product")
		}
		if existingProduct != nil && existingProduct.ID != id {
			return nil, errors.New("product with similar name already exists")
		}

		product.Name = req.Name
		product.Slug = newSlug
	}

	// Update description if provided
	if req.Description != "" {
		product.Description = req.Description
	}

	// Update price if provided
	if req.Price > 0 {
		product.Price = req.Price
	}

	// Update compare price if provided
	if req.ComparePrice >= 0 {
		product.ComparePrice = req.ComparePrice
	}

	// Update SKU if provided
	if req.SKU != "" && req.SKU != product.SKU {
		exists, err := s.productRepo.SKUExists(req.SKU)
		if err != nil {
			logger.Error("Failed to check SKU existence: %v", err)
			return nil, errors.New("failed to update product")
		}
		if exists {
			return nil, errors.New("SKU already exists")
		}
		product.SKU = req.SKU
	}

	// Update barcode if provided
	if req.Barcode != "" {
		product.Barcode = req.Barcode
	}

	// Update quantity if provided
	if req.Quantity >= 0 {
		product.Quantity = req.Quantity
	}

	// Update images if provided
	if req.Images != nil {
		imagesBytes, _ := json.Marshal(req.Images)
		product.Images = datatypes.JSON(imagesBytes)
	}

	// Update category if provided
	if req.CategoryID > 0 && req.CategoryID != product.CategoryID {
		category, err := s.categoryRepo.GetByID(req.CategoryID)
		if err != nil {
			logger.Error("Failed to check category: %v", err)
			return nil, errors.New("failed to update product")
		}
		if category == nil {
			return nil, errors.New("category not found")
		}
		product.CategoryID = req.CategoryID
	}

	// Update tags if provided
	if req.TagIDs != nil {
		var tags []models.Tag
		for _, tagID := range req.TagIDs {
			tag, err := s.tagRepo.GetByID(tagID)
			if err != nil {
				logger.Error("Failed to check tag %d: %v", tagID, err)
				return nil, errors.New("failed to update product")
			}
			if tag == nil {
				return nil, fmt.Errorf("tag with ID %d not found", tagID)
			}
			tags = append(tags, *tag)
		}

		// Update tags association
		if err := s.productRepo.UpdateTags(product, tags); err != nil {
			logger.Error("Failed to update product tags: %v", err)
			return nil, errors.New("failed to update product")
		}
	}

	// Update is_active if provided
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	// Update is_featured if provided
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}

	// Save changes
	if err := s.productRepo.Update(product); err != nil {
		logger.Error("Failed to update product %d: %v", id, err)
		return nil, errors.New("failed to update product")
	}

	logger.Success("Product updated successfully: %s (id: %d)", product.Name, product.ID)

	// Reload product with relationships
	product, _ = s.productRepo.GetByID(id)
	response := product.ToResponse()
	return &response, nil
}

// ========================================
// DELETE PRODUCT
// ========================================

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(id uint) error {
	// Get product
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get product %d: %v", id, err)
		return errors.New("failed to delete product")
	}
	if product == nil {
		return errors.New("product not found")
	}

	// TODO: Check if product is in any orders (when order model is ready)

	// Delete product
	if err := s.productRepo.Delete(id); err != nil {
		logger.Error("Failed to delete product %d: %v", id, err)
		return errors.New("failed to delete product")
	}

	logger.Success("Product deleted successfully: %s (id: %d)", product.Name, id)
	return nil
}

// ========================================
// STOCK MANAGEMENT
// ========================================

// UpdateStock updates product stock quantity
func (s *ProductService) UpdateStock(productID uint, quantity int) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	// Check product exists
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		logger.Error("Failed to get product %d: %v", productID, err)
		return errors.New("failed to update stock")
	}
	if product == nil {
		return errors.New("product not found")
	}

	if err := s.productRepo.UpdateStock(productID, quantity); err != nil {
		logger.Error("Failed to update stock for product %d: %v", productID, err)
		return errors.New("failed to update stock")
	}

	logger.Success("Stock updated for product %d: %d units", productID, quantity)
	return nil
}
