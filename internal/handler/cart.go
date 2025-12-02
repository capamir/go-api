package handler

import (
	"strconv"

	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// CartHandler handles cart HTTP requests
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler creates a new cart handler
func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// ========================================
// AUTHENTICATED ENDPOINTS (Require JWT)
// ========================================

// GetUserCart retrieves user's cart
// @Summary Get user's cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.CartResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /cart [get]
func (h *CartHandler) GetUserCart(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Get cart
	cart, err := h.cartService.GetUserCart(userID)
	if err != nil {
		logger.Error("Failed to get cart for user %d: %v", userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Cart retrieved successfully", cart)
}

// AddToCart adds product to cart
// @Summary Add product to cart
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.AddToCartRequest true "Cart item data"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /cart/items [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	var req service.AddToCartRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid add to cart request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Add to cart
	cart, err := h.cartService.AddToCart(userID, req)
	if err != nil {
		logger.Error("Failed to add to cart for user %d: %v", userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Product added to cart successfully", cart)
}

// UpdateCartItem updates cart item quantity
// @Summary Update cart item quantity
// @Tags cart
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param item_id path int true "Cart item ID"
// @Param request body service.UpdateCartItemRequest true "Updated quantity"
// @Success 200 {object} models.CartResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /cart/items/{item_id} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Parse item ID from URL
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid cart item ID: %v", err)
		utils.RespondBadRequest(c, "Invalid cart item ID")
		return
	}

	var req service.UpdateCartItemRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update cart item request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Update cart item
	cart, err := h.cartService.UpdateCartItem(userID, uint(itemID), req)
	if err != nil {
		logger.Error("Failed to update cart item %d for user %d: %v", itemID, userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Cart item updated successfully", cart)
}

// RemoveCartItem removes item from cart
// @Summary Remove item from cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Param item_id path int true "Cart item ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /cart/items/{item_id} [delete]
func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Parse item ID from URL
	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid cart item ID: %v", err)
		utils.RespondBadRequest(c, "Invalid cart item ID")
		return
	}

	// Remove cart item
	if err := h.cartService.RemoveCartItem(userID, uint(itemID)); err != nil {
		logger.Error("Failed to remove cart item %d for user %d: %v", itemID, userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Cart item removed successfully", nil)
}

// ClearCart clears all items from cart
// @Summary Clear entire cart
// @Tags cart
// @Security BearerAuth
// @Produce json
// @Success 200 {object} utils.SuccessResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Clear cart
	if err := h.cartService.ClearCart(userID); err != nil {
		logger.Error("Failed to clear cart for user %d: %v", userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Cart cleared successfully", nil)
}
