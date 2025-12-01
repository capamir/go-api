package service

import (
	"errors"
	"fmt"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
)

// CategoryService handles category business logic
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// ========================================
// DTOs (Data Transfer Objects)
// ========================================

// CreateCategoryRequest represents category creation data
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// UpdateCategoryRequest represents category update data
type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	ParentID    *uint  `json:"parent_id" binding:"omitempty"`
	ImageURL    string `json:"image_url" binding:"omitempty,url"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

// CategoryListResponse represents paginated category list
type CategoryListResponse struct {
	Categories []models.CategoryResponse `json:"categories"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

// ========================================
// CREATE CATEGORY
// ========================================

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(req CreateCategoryRequest) (*models.CategoryResponse, error) {
	// Generate slug from name
	slug := utils.GenerateSlug(req.Name)

	// Check if slug already exists
	exists, err := s.categoryRepo.SlugExists(slug)
	if err != nil {
		logger.Error("Failed to check slug existence: %v", err)
		return nil, errors.New("failed to create category")
	}
	if exists {
		return nil, errors.New("category with similar name already exists")
	}

	// Validate parent category if provided
	if req.ParentID != nil {
		parent, err := s.categoryRepo.GetByID(*req.ParentID)
		if err != nil {
			logger.Error("Failed to check parent category: %v", err)
			return nil, errors.New("failed to create category")
		}
		if parent == nil {
			return nil, errors.New("parent category not found")
		}
	}

	// Set default for IsActive
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create category
	category := &models.Category{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		ImageURL:    req.ImageURL,
		IsActive:    isActive,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		logger.Error("Failed to create category: %v", err)
		return nil, errors.New("failed to create category")
	}

	logger.Success("Category created successfully: %s (slug: %s)", category.Name, category.Slug)

	response := category.ToResponse()
	return &response, nil
}

// ========================================
// GET CATEGORY
// ========================================

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*models.CategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get category %d: %v", id, err)
		return nil, errors.New("failed to get category")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	response := category.ToResponse()
	return &response, nil
}

// GetCategoryBySlug retrieves a category by slug
func (s *CategoryService) GetCategoryBySlug(slug string) (*models.CategoryResponse, error) {
	category, err := s.categoryRepo.GetBySlug(slug)
	if err != nil {
		logger.Error("Failed to get category by slug %s: %v", slug, err)
		return nil, errors.New("failed to get category")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	response := category.ToResponse()
	return &response, nil
}

// GetCategoryWithChildren retrieves a category with its subcategories
func (s *CategoryService) GetCategoryWithChildren(id uint) (*models.CategoryWithChildren, error) {
	category, err := s.categoryRepo.GetWithChildren(id)
	if err != nil {
		logger.Error("Failed to get category with children %d: %v", id, err)
		return nil, errors.New("failed to get category")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	response := category.ToResponseWithChildren()
	return &response, nil
}

// ========================================
// LIST CATEGORIES
// ========================================

// GetAllCategories retrieves all categories with pagination
func (s *CategoryService) GetAllCategories(page, limit int) (*CategoryListResponse, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	categories, total, err := s.categoryRepo.GetAll(limit, offset)
	if err != nil {
		logger.Error("Failed to get categories: %v", err)
		return nil, errors.New("failed to get categories")
	}

	// Convert to response DTOs
	categoryResponses := make([]models.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = category.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &CategoryListResponse{
		Categories: categoryResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetActiveCategories retrieves only active categories
func (s *CategoryService) GetActiveCategories(page, limit int) (*CategoryListResponse, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	categories, total, err := s.categoryRepo.GetActiveCategories(limit, offset)
	if err != nil {
		logger.Error("Failed to get active categories: %v", err)
		return nil, errors.New("failed to get categories")
	}

	// Convert to response DTOs
	categoryResponses := make([]models.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = category.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &CategoryListResponse{
		Categories: categoryResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetRootCategories retrieves top-level categories (no parent)
func (s *CategoryService) GetRootCategories() ([]models.CategoryResponse, error) {
	categories, err := s.categoryRepo.GetRootCategories()
	if err != nil {
		logger.Error("Failed to get root categories: %v", err)
		return nil, errors.New("failed to get categories")
	}

	// Convert to response DTOs
	categoryResponses := make([]models.CategoryResponse, len(categories))
	for i, category := range categories {
		categoryResponses[i] = category.ToResponse()
	}

	return categoryResponses, nil
}

// GetCategoryTree retrieves all categories in hierarchical structure
func (s *CategoryService) GetCategoryTree() ([]models.CategoryWithChildren, error) {
	categories, err := s.categoryRepo.GetCategoryTree()
	if err != nil {
		logger.Error("Failed to get category tree: %v", err)
		return nil, errors.New("failed to get category tree")
	}

	// Convert to response DTOs with children
	categoryResponses := make([]models.CategoryWithChildren, len(categories))
	for i, category := range categories {
		categoryResponses[i] = category.ToResponseWithChildren()
	}

	return categoryResponses, nil
}

// ========================================
// UPDATE CATEGORY
// ========================================

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id uint, req UpdateCategoryRequest) (*models.CategoryResponse, error) {
	// Get existing category
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get category %d: %v", id, err)
		return nil, errors.New("failed to update category")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Update name and slug if provided
	if req.Name != "" && req.Name != category.Name {
		newSlug := utils.GenerateSlug(req.Name)

		// Check if new slug already exists (excluding current category)
		existingCategory, err := s.categoryRepo.GetBySlug(newSlug)
		if err != nil {
			logger.Error("Failed to check slug existence: %v", err)
			return nil, errors.New("failed to update category")
		}
		if existingCategory != nil && existingCategory.ID != id {
			return nil, errors.New("category with similar name already exists")
		}

		category.Name = req.Name
		category.Slug = newSlug
	}

	// Update description if provided
	if req.Description != "" {
		category.Description = req.Description
	}

	// Update parent if provided
	if req.ParentID != nil {
		// Prevent self-reference
		if *req.ParentID == id {
			return nil, errors.New("category cannot be its own parent")
		}

		// Validate parent exists
		if *req.ParentID != 0 {
			parent, err := s.categoryRepo.GetByID(*req.ParentID)
			if err != nil {
				logger.Error("Failed to check parent category: %v", err)
				return nil, errors.New("failed to update category")
			}
			if parent == nil {
				return nil, errors.New("parent category not found")
			}
		}

		category.ParentID = req.ParentID
	}

	// Update image URL if provided
	if req.ImageURL != "" {
		category.ImageURL = req.ImageURL
	}

	// Update is_active if provided
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	// Save changes
	if err := s.categoryRepo.Update(category); err != nil {
		logger.Error("Failed to update category %d: %v", id, err)
		return nil, errors.New("failed to update category")
	}

	logger.Success("Category updated successfully: %s (id: %d)", category.Name, category.ID)

	response := category.ToResponse()
	return &response, nil
}

// ========================================
// DELETE CATEGORY
// ========================================

// DeleteCategory soft deletes a category
func (s *CategoryService) DeleteCategory(id uint) error {
	// Get category
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get category %d: %v", id, err)
		return errors.New("failed to delete category")
	}
	if category == nil {
		return errors.New("category not found")
	}

	// Check if category has children
	categoryWithChildren, err := s.categoryRepo.GetWithChildren(id)
	if err != nil {
		logger.Error("Failed to check category children: %v", err)
		return errors.New("failed to delete category")
	}
	if len(categoryWithChildren.Children) > 0 {
		return fmt.Errorf("cannot delete category with %d subcategories. Delete or reassign them first", len(categoryWithChildren.Children))
	}

	// TODO: Check if category has products (when product model is ready)
	// For now, just delete the category

	// Delete category
	if err := s.categoryRepo.Delete(id); err != nil {
		logger.Error("Failed to delete category %d: %v", id, err)
		return errors.New("failed to delete category")
	}

	logger.Success("Category deleted successfully: %s (id: %d)", category.Name, id)
	return nil
}
