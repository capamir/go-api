package handler

import (
	"strconv"

	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

// OrderHandler handles order HTTP requests
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// ========================================
// USER ENDPOINTS (Authenticated)
// ========================================

// GetUserOrders retrieves user's order history
// @Summary Get user's orders
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} service.OrderListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Get user orders
	orders, err := h.orderService.GetUserOrders(userID, page, limit)
	if err != nil {
		logger.Error("Failed to get orders for user %d: %v", userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Orders retrieved successfully", orders)
}

// GetOrderByID retrieves a single order
// @Summary Get user's order by ID
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.OrderResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Parse order ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid order ID: %v", err)
		utils.RespondBadRequest(c, "Invalid order ID")
		return
	}

	// Get order
	order, err := h.orderService.GetOrderByID(userID, uint(id))
	if err != nil {
		logger.Error("Failed to get order %d for user %d: %v", id, userID, err)
		if err.Error() == "order not found" {
			utils.RespondNotFound(c, err.Error())
		} else {
			utils.RespondBadRequest(c, err.Error())
		}
		return
	}

	utils.RespondSuccess(c, "Order retrieved successfully", order)
}

// CancelOrder cancels an order
// @Summary Cancel user's order
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	// Parse order ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid order ID: %v", err)
		utils.RespondBadRequest(c, "Invalid order ID")
		return
	}

	// Cancel order
	if err := h.orderService.CancelOrder(userID, uint(id)); err != nil {
		logger.Error("Failed to cancel order %d for user %d: %v", id, userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Order cancelled successfully", nil)
}

// CreateOrder creates order from cart (checkout)
// @Summary Create order (checkout)
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body service.CreateOrderRequest true "Checkout data"
// @Success 201 {object} models.OrderResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	// Get user ID from JWT
	userID, exists := middleware.GetUserID(c)
	if !exists {
		logger.Warn("User ID not found in context")
		utils.RespondUnauthorized(c, "User not authenticated")
		return
	}

	var req service.CreateOrderRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid create order request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Create order
	order, err := h.orderService.CreateOrder(userID, req)
	if err != nil {
		logger.Error("Failed to create order for user %d: %v", userID, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondCreated(c, "Order created successfully", order)
}

// ========================================
// ADMIN ENDPOINTS (Require Admin Role)
// ========================================

// GetAllOrders retrieves all orders (admin only)
// @Summary Get all orders (admin)
// @Tags admin/orders
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status"
// @Param payment_status query string false "Filter by payment status"
// @Param user_id query int false "Filter by user ID"
// @Param min_total query float64 false "Minimum total"
// @Param max_total query float64 false "Maximum total"
// @Param q query string false "Search by order number"
// @Success 200 {object} service.OrderListResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /admin/orders [get]
func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	// Check admin role
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to get all orders")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	var filterReq service.OrderFilterRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&filterReq); err != nil {
		logger.Warn("Invalid order filter parameters: %v", err)
		utils.RespondBadRequest(c, "Invalid query parameters")
		return
	}

	// Get all orders
	orders, err := h.orderService.GetAllOrders(filterReq)
	if err != nil {
		logger.Error("Failed to get all orders: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Orders retrieved successfully", orders)
}

// GetOrderStats retrieves order statistics (admin only)
// @Summary Get order statistics (admin)
// @Tags admin/orders
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Router /admin/orders/stats [get]
func (h *OrderHandler) GetOrderStats(c *gin.Context) {
	// Check admin role
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to get order stats")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Get statistics
	stats, err := h.orderService.GetOrderStats()
	if err != nil {
		logger.Error("Failed to get order stats: %v", err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Order statistics retrieved successfully", stats)
}

// UpdateOrderStatus updates order status (admin only)
// @Summary Update order status (admin)
// @Tags admin/orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param request body service.UpdateOrderStatusRequest true "New status"
// @Success 200 {object} models.OrderResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 403 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /admin/orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	// Check admin role
	userRole, exists := middleware.GetUserRole(c)
	if !exists || userRole != "admin" {
		logger.Warn("Non-admin user attempted to update order status")
		utils.RespondForbidden(c, "Admin access required")
		return
	}

	// Parse order ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		logger.Warn("Invalid order ID: %v", err)
		utils.RespondBadRequest(c, "Invalid order ID")
		return
	}

	var req service.UpdateOrderStatusRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid update order status request: %v", err)
		utils.RespondBadRequest(c, "Invalid request data")
		return
	}

	// Update order status
	order, err := h.orderService.UpdateOrderStatus(uint(id), req)
	if err != nil {
		logger.Error("Failed to update order status %d: %v", id, err)
		utils.RespondBadRequest(c, err.Error())
		return
	}

	utils.RespondSuccess(c, "Order status updated successfully", order)
}
