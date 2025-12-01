package handler

import (
	"strconv"

	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// ========================================
// PUBLIC ENDPOINTS
// ========================================

// GetAllProducts lists all products with filters and pagination
// @Summary Get all products
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param category_id query int false "Filter by category ID"
// @Param tag_ids query []int false "Filter by tag IDs"
// @Param min_price query float64 false "Minimum price"
// @Param max_price query float64 false "Maximum price"
// @Param in_stock query bool false "Only in stock products"
// @Param is_active query bool false "Only active products"
// @Param is_featured query bool false "Only featured products"
// @Param q query string false "Search query"
// @Success 200 {object} service.ProductListResponse
// @Router /products [get]
func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	var filterReq service.ProductFilterRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&filterReq); err != nil {
		logger.Warn("Invalid query parameters: %v", err)
		utils.RespondBadRequest(c, "Invalid query parameters")
		return
	}

	// Get products
	response, err := h.productService.GetAllProducts(filterReq)
	if err != nil {
		logger.Error("Failed to get products: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Products retrieved successfully", response)
}

// GetFeaturedProducts retrieves featured products
// @Summary Get featured products
// @Tags products
// @Produce json
// @Param limit query int false "Number of products" default(10)
// @Success 200 {array} models.ProductResponse
// @Router /products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		logger.Error("Failed to get featured products: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Featured products retrieved successfully", products)
}

// GetProductByID retrieves a single product by ID
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.ProductResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid product ID: %v", err)
		utils.RespondBadRequest(c, "Invalid product ID")
		return
	}

	// Get product
	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		logger.Error("Failed to get product %d: %v", id, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Product retrieved successfully", product)
}

// GetProductBySlug retrieves a single product by slug
// @Summary Get product by slug
// @Tags products
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} models.ProductResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.RespondBadRequest(c, "Slug is required")
		return
	}

	// Get product
	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		logger.Error("Failed to get product by slug %s: %v", slug, err)
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Product retrieved successfully", product)
}

// GetRelatedProducts retrieves related products (same category)
// @Summary Get related products
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Param limit query int false "Number of products" default(4)
// @Success 200 {array} models.ProductResponse
// @Router /products/{id}/related [get]
func (h *ProductHandler) GetRelatedProducts(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid product ID: %v", err)
		utils.RespondBadRequest(c, "Invalid product ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))

	// Get related products
	products, err := h.productService.GetRelatedProducts(uint(id), limit)
	if err != nil {
		logger.Error("Failed to get related products for %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Related products retrieved successfully", products)
}

// ========================================
// ADMIN ENDPOINTS (Protected)
// ========================================

// CreateProduct creates a new product (Admin only)
// @Summary Create product
// @Tags admin/products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateProductRequest true "Product data"
// @Success 201 {object} models.ProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /admin/products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to create product")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	var req service.CreateProductRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid create product request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Create product
	product, err := h.productService.CreateProduct(req)
	if err != nil {
		logger.Error("Failed to create product: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondCreated(c, "Product created successfully", product)
}

// UpdateProduct updates an existing product (Admin only)
// @Summary Update product
// @Tags admin/products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body service.UpdateProductRequest true "Product data"
// @Success 200 {object} models.ProductResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to update product")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid product ID: %v", err)
		utils.RespondBadRequest(c, "Invalid product ID")
		return
	}

	var req service.UpdateProductRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update product request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Update product
	product, err := h.productService.UpdateProduct(uint(id), req)
	if err != nil {
		logger.Error("Failed to update product %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Product updated successfully", product)
}

// DeleteProduct deletes a product (Admin only)
// @Summary Delete product
// @Tags admin/products
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to delete product")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid product ID: %v", err)
		utils.RespondBadRequest(c, "Invalid product ID")
		return
	}

	// Delete product
	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		logger.Error("Failed to delete product %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Product deleted successfully", nil)
}

// UpdateStock updates product stock (Admin only)
// @Summary Update product stock
// @Tags admin/products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body object true "Stock data" example({"quantity": 100})
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /admin/products/{id}/stock [put]
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	// Get user from context
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to update stock")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid product ID: %v", err)
		utils.RespondBadRequest(c, "Invalid product ID")
		return
	}

	// Parse request body
	var req struct {
		Quantity int `json:"quantity" binding:"required,gte=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid stock update request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Update stock
	if err := h.productService.UpdateStock(uint(id), req.Quantity); err != nil {
		logger.Error("Failed to update stock for product %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Stock updated successfully", gin.H{
		"product_id": id,
		"quantity":   req.Quantity,
	})
}
