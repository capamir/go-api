package service

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/pkg/logger"
	"gorm.io/datatypes"
)

// OrderService handles order business logic
type OrderService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

// NewOrderService creates a new order service
func NewOrderService(
	orderRepo *repository.OrderRepository,
	cartRepo *repository.CartRepository,
	productRepo *repository.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// ========================================
// DTOs (Data Transfer Objects)
// ========================================

// CreateOrderRequest represents checkout data
type CreateOrderRequest struct {
	ShippingAddress models.ShippingAddressData `json:"shipping_address" binding:"required"`
	Notes           string                     `json:"notes" binding:"omitempty,max=500"`
}

// UpdateOrderStatusRequest represents status update
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing shipped delivered cancelled"`
}

// OrderFilterRequest represents order filter parameters
type OrderFilterRequest struct {
	Status        string  `form:"status"`
	PaymentStatus string  `form:"payment_status"`
	UserID        *uint   `form:"user_id"`
	MinTotal      *float64 `form:"min_total"`
	MaxTotal      *float64 `form:"max_total"`
	SearchQuery   string  `form:"q"`
	Page          int     `form:"page"`
	Limit         int     `form:"limit"`
}

// OrderListResponse represents paginated order list
type OrderListResponse struct {
	Orders     []models.OrderResponse `json:"orders"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}

// ========================================
// CREATE ORDER (CHECKOUT)
// ========================================

// CreateOrder creates an order from user's cart
func (s *OrderService) CreateOrder(userID uint, req CreateOrderRequest) (*models.OrderResponse, error) {
	// Get and validate cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		logger.Error("Failed to get cart for user %d: %v", userID, err)
		return nil, errors.New("failed to create order")
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate stock and calculate totals
	var subtotal float64
	var orderItems []models.OrderItem

	for _, cartItem := range cart.Items {
		// Validate product
		product, err := s.productRepo.GetByID(cartItem.ProductID)
		if err != nil {
			logger.Error("Failed to get product %d: %v", cartItem.ProductID, err)
			return nil, errors.New("failed to create order")
		}
		if product == nil {
			return nil, fmt.Errorf("product not found: %d", cartItem.ProductID)
		}

		// Check if product is active
		if !product.IsActive {
			return nil, fmt.Errorf("product '%s' is no longer available", product.Name)
		}

		// Check stock
		if product.Quantity < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for '%s'. Available: %d, requested: %d",
				product.Name, product.Quantity, cartItem.Quantity)
		}

		// Create order item (snapshot)
		itemSubtotal := cartItem.Price * float64(cartItem.Quantity)
		orderItem := models.OrderItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			ProductSKU:  product.SKU,
			Quantity:    cartItem.Quantity,
			Price:       cartItem.Price, // Use cart price (may differ from current)
			Subtotal:    itemSubtotal,
		}

		orderItems = append(orderItems, orderItem)
		subtotal += itemSubtotal
	}

	// Calculate tax and shipping (can be customized)
	taxRate := 0.10 // 10% tax
	tax := subtotal * taxRate
	shippingCost := 0.0 // Free shipping for now

	if subtotal < 100 {
		shippingCost = 10.0 // $10 shipping for orders under $100
	}

	total := subtotal + tax + shippingCost

	// Convert shipping address to JSON
	addressBytes, _ := json.Marshal(req.ShippingAddress)
	shippingAddressJSON := datatypes.JSON(addressBytes)

	// Generate order number
	orderNumber, err := s.orderRepo.GenerateOrderNumber()
	if err != nil {
		logger.Error("Failed to generate order number: %v", err)
		return nil, errors.New("failed to create order")
	}

	// Create order
	order := &models.Order{
		OrderNumber:     orderNumber,
		UserID:          userID,
		Items:           orderItems,
		Subtotal:        subtotal,
		Tax:             tax,
		ShippingCost:    shippingCost,
		Total:           total,
		Status:          models.OrderStatusPending,
		PaymentStatus:   models.PaymentStatusPending,
		ShippingAddress: shippingAddressJSON,
		Notes:           req.Notes,
	}

	if err := s.orderRepo.Create(order); err != nil {
		logger.Error("Failed to create order: %v", err)
		return nil, errors.New("failed to create order")
	}

	// Decrement product stock
	for _, item := range orderItems {
		if err := s.productRepo.DecrementStock(item.ProductID, item.Quantity); err != nil {
			logger.Error("Failed to decrement stock for product %d: %v", item.ProductID, err)
			// Note: In production, this should rollback the entire transaction
		}
	}

	// Clear cart
	if err := s.cartRepo.ClearCart(cart.ID); err != nil {
		logger.Warn("Failed to clear cart after order: %v", err)
		// Non-critical, continue
	}

	logger.Success("Order created successfully: %s for user %d (Total: $%.2f)", orderNumber, userID, total)

	// Reload order with relations
	createdOrder, _ := s.orderRepo.GetByID(order.ID)
	response := createdOrder.ToResponse()
	return &response, nil
}

// ========================================
// GET ORDERS
// ========================================

// GetUserOrders retrieves user's order history
func (s *OrderService) GetUserOrders(userID uint, page, limit int) (*OrderListResponse, error) {
	// Default pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	orders, total, err := s.orderRepo.GetUserOrders(userID, limit, offset)
	if err != nil {
		logger.Error("Failed to get orders for user %d: %v", userID, err)
		return nil, errors.New("failed to get orders")
	}

	// Convert to response DTOs
	orderResponses := make([]models.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = order.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &OrderListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// GetOrderByID retrieves a single order
func (s *OrderService) GetOrderByID(userID, orderID uint) (*models.OrderResponse, error) {
	// Check ownership
	isOwner, err := s.orderRepo.IsOrderOwner(orderID, userID)
	if err != nil {
		logger.Error("Failed to check order ownership: %v", err)
		return nil, errors.New("failed to get order")
	}
	if !isOwner {
		return nil, errors.New("order not found")
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		logger.Error("Failed to get order %d: %v", orderID, err)
		return nil, errors.New("failed to get order")
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	response := order.ToResponse()
	return &response, nil
}

// ========================================
// CANCEL ORDER
// ========================================

// CancelOrder cancels an order (user)
func (s *OrderService) CancelOrder(userID, orderID uint) error {
	// Check ownership
	isOwner, err := s.orderRepo.IsOrderOwner(orderID, userID)
	if err != nil {
		logger.Error("Failed to check order ownership: %v", err)
		return errors.New("failed to cancel order")
	}
	if !isOwner {
		return errors.New("order not found")
	}

	// Get order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		logger.Error("Failed to get order %d: %v", orderID, err)
		return errors.New("failed to cancel order")
	}
	if order == nil {
		return errors.New("order not found")
	}

	// Check if can be cancelled
	if !order.CanBeCancelled() {
		return fmt.Errorf("order cannot be cancelled (status: %s)", order.Status)
	}

	// Update status
	order.Status = models.OrderStatusCancelled
	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("Failed to update order %d: %v", orderID, err)
		return errors.New("failed to cancel order")
	}

	// Restore stock
	for _, item := range order.Items {
		if err := s.productRepo.UpdateStock(item.ProductID, item.Product.Quantity+item.Quantity); err != nil {
			logger.Warn("Failed to restore stock for product %d: %v", item.ProductID, err)
		}
	}

	logger.Success("Order %s cancelled by user %d", order.OrderNumber, userID)
	return nil
}

// ========================================
// ADMIN OPERATIONS
// ========================================

// GetAllOrders retrieves all orders (admin)
func (s *OrderService) GetAllOrders(filterReq OrderFilterRequest) (*OrderListResponse, error) {
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
	filters := &repository.OrderFilterOptions{
		Status:        filterReq.Status,
		PaymentStatus: filterReq.PaymentStatus,
		UserID:        filterReq.UserID,
		MinTotal:      filterReq.MinTotal,
		MaxTotal:      filterReq.MaxTotal,
		SearchQuery:   filterReq.SearchQuery,
	}

	orders, total, err := s.orderRepo.GetAll(limit, offset, filters)
	if err != nil {
		logger.Error("Failed to get all orders: %v", err)
		return nil, errors.New("failed to get orders")
	}

	// Convert to response DTOs
	orderResponses := make([]models.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = order.ToResponse()
	}

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &OrderListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatus updates order status (admin)
func (s *OrderService) UpdateOrderStatus(orderID uint, req UpdateOrderStatusRequest) (*models.OrderResponse, error) {
	// Get order
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		logger.Error("Failed to get order %d: %v", orderID, err)
		return nil, errors.New("failed to update order")
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	// Validate status transition
	if !order.CanBeUpdated() {
		return nil, fmt.Errorf("order cannot be updated (status: %s)", order.Status)
	}

	// Update status
	order.Status = req.Status
	if err := s.orderRepo.Update(order); err != nil {
		logger.Error("Failed to update order %d: %v", orderID, err)
		return nil, errors.New("failed to update order")
	}

	logger.Success("Order %s status updated to '%s'", order.OrderNumber, req.Status)

	response := order.ToResponse()
	return &response, nil
}

// GetOrderStats retrieves order statistics (admin)
func (s *OrderService) GetOrderStats() (map[string]interface{}, error) {
	stats, err := s.orderRepo.GetOrderStats()
	if err != nil {
		logger.Error("Failed to get order stats: %v", err)
		return nil, errors.New("failed to get statistics")
	}
	return stats, nil
}
