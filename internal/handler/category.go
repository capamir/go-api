package handler

import (
	"strconv"

	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// ========================================
// PUBLIC ENDPOINTS
// ========================================

// GetAllCategories lists all categories with pagination
// @Summary Get all categories
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} service.CategoryListResponse
// @Router /categories [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Get categories
	response, err := h.categoryService.GetAllCategories(page, limit)
	if err != nil {
		logger.Error("Failed to get categories: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Categories retrieved successfully", response)
}

// GetActiveCategories lists only active categories
// @Summary Get active categories
// @Tags categories
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} service.CategoryListResponse
// @Router /categories/active [get]
func (h *CategoryHandler) GetActiveCategories(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Get active categories
	response, err := h.categoryService.GetActiveCategories(page, limit)
	if err != nil {
		logger.Error("Failed to get active categories: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Active categories retrieved successfully", response)
}

// GetRootCategories lists top-level categories
// @Summary Get root categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.CategoryResponse
// @Router /categories/root [get]
func (h *CategoryHandler) GetRootCategories(c *gin.Context) {
	categories, err := h.categoryService.GetRootCategories()
	if err != nil {
		logger.Error("Failed to get root categories: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Root categories retrieved successfully", categories)
}

// GetCategoryTree retrieves full category hierarchy
// @Summary Get category tree
// @Tags categories
// @Produce json
// @Success 200 {array} models.CategoryWithChildren
// @Router /categories/tree [get]
func (h *CategoryHandler) GetCategoryTree(c *gin.Context) {
	tree, err := h.categoryService.GetCategoryTree()
	if err != nil {
		logger.Error("Failed to get category tree: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Category tree retrieved successfully", tree)
}

// GetCategoryByID retrieves a single category by ID
// @Summary Get category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.CategoryResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid category ID: %v", err)
		utils.RespondBadRequest(c, "Invalid category ID")
		return
	}

	// Get category
	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		logger.Error("Failed to get category %d: %v", id, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Category retrieved successfully", category)
}

// GetCategoryBySlug retrieves a single category by slug
// @Summary Get category by slug
// @Tags categories
// @Produce json
// @Param slug path string true "Category slug"
// @Success 200 {object} models.CategoryResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /categories/slug/{slug} [get]
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.RespondBadRequest(c, "Slug is required")
		return
	}

	// Get category
	category, err := h.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		logger.Error("Failed to get category by slug %s: %v", slug, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Category retrieved successfully", category)
}

// GetCategoryWithChildren retrieves a category with its subcategories
// @Summary Get category with children
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.CategoryWithChildren
// @Failure 404 {object} utils.ErrorResponse
// @Router /categories/{id}/children [get]
func (h *CategoryHandler) GetCategoryWithChildren(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid category ID: %v", err)
		utils.RespondBadRequest(c, "Invalid category ID")
		return
	}

	// Get category with children
	category, err := h.categoryService.GetCategoryWithChildren(uint(id))
	if err != nil {
		logger.Error("Failed to get category with children %d: %v", id, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Category with children retrieved successfully", category)
}

// ========================================
// ADMIN ENDPOINTS (Protected)
// ========================================

// CreateCategory creates a new category (Admin only)
// @Summary Create category
// @Tags admin/categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.CategoryResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /admin/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to create category")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	var req service.CreateCategoryRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid create category request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Create category
	category, err := h.categoryService.CreateCategory(req)
	if err != nil {
		logger.Error("Failed to create category: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondCreated(c, "Category created successfully", category)
}

// UpdateCategory updates an existing category (Admin only)
// @Summary Update category
// @Tags admin/categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param request body service.UpdateCategoryRequest true "Category data"
// @Success 200 {object} models.CategoryResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to update category")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid category ID: %v", err)
		utils.RespondBadRequest(c, "Invalid category ID")
		return
	}

	var req service.UpdateCategoryRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update category request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Update category
	category, err := h.categoryService.UpdateCategory(uint(id), req)
	if err != nil {
		logger.Error("Failed to update category %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Category updated successfully", category)
}

// DeleteCategory deletes a category (Admin only)
// @Summary Delete category
// @Tags admin/categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to delete category")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid category ID: %v", err)
		utils.RespondBadRequest(c, "Invalid category ID")
		return
	}

	// Delete category
	if err := h.categoryService.DeleteCategory(uint(id)); err != nil {
		logger.Error("Failed to delete category %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Category deleted successfully", nil)
}
