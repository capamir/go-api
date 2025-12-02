package service

import (
	"errors"
	"fmt"

	"github.com/capamir/go-api/internal/models"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/pkg/logger"
)

// CartService handles cart business logic
type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

// NewCartService creates a new cart service
func NewCartService(
	cartRepo *repository.CartRepository,
	productRepo *repository.ProductRepository,
) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// ========================================
// DTOs (Data Transfer Objects)
// ========================================

// AddToCartRequest represents adding item to cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest represents updating cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// ========================================
// GET CART
// ========================================

// GetUserCart retrieves user's cart with all items
func (s *CartService) GetUserCart(userID uint) (*models.CartResponse, error) {
	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		logger.Error("Failed to get cart for user %d: %v", userID, err)
		return nil, errors.New("failed to get cart")
	}

	response := cart.ToResponse()
	return &response, nil
}

// ========================================
// ADD TO CART
// ========================================

// AddToCart adds a product to user's cart
func (s *CartService) AddToCart(userID uint, req AddToCartRequest) (*models.CartResponse, error) {
	// Validate product exists and is available
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		logger.Error("Failed to get product %d: %v", req.ProductID, err)
		return nil, errors.New("failed to add to cart")
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Check if product is active
	if !product.IsActive {
		return nil, errors.New("product is not available")
	}

	// Check stock availability
	if product.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock. Available: %d", product.Quantity)
	}

	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(userID)
	if err != nil {
		logger.Error("Failed to get cart for user %d: %v", userID, err)
		return nil, errors.New("failed to add to cart")
	}

	// Check if product already in cart
	existingItem, err := s.cartRepo.GetCartItemByProductID(cart.ID, req.ProductID)
	if err != nil {
		logger.Error("Failed to check existing cart item: %v", err)
		return nil, errors.New("failed to add to cart")
	}

	if existingItem != nil {
		// Update quantity of existing item
		newQuantity := existingItem.Quantity + req.Quantity

		// Check stock for new quantity
		if product.Quantity < newQuantity {
			return nil, fmt.Errorf("insufficient stock. Available: %d, requested: %d", product.Quantity, newQuantity)
		}

		existingItem.Quantity = newQuantity
		if err := s.cartRepo.UpdateItem(existingItem); err != nil {
			logger.Error("Failed to update cart item: %v", err)
			return nil, errors.New("failed to add to cart")
		}

		logger.Success("Updated cart item quantity for user %d, product %d", userID, req.ProductID)
	} else {
		// Add new item to cart
		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			Price:     product.Price, // Snapshot current price
		}

		if err := s.cartRepo.AddItem(cartItem); err != nil {
			logger.Error("Failed to add cart item: %v", err)
			return nil, errors.New("failed to add to cart")
		}

		logger.Success("Added product %d to cart for user %d", req.ProductID, userID)
	}

	// Reload cart with updated items
	updatedCart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		logger.Error("Failed to reload cart: %v", err)
		return nil, errors.New("failed to get updated cart")
	}

	response := updatedCart.ToResponse()
	return &response, nil
}

// ========================================
// UPDATE CART ITEM
// ========================================

// UpdateCartItem updates cart item quantity
func (s *CartService) UpdateCartItem(userID, itemID uint, req UpdateCartItemRequest) (*models.CartResponse, error) {
	// Check ownership
	isOwner, err := s.cartRepo.IsCartItemOwner(itemID, userID)
	if err != nil {
		logger.Error("Failed to check cart item ownership: %v", err)
		return nil, errors.New("failed to update cart item")
	}
	if !isOwner {
		return nil, errors.New("cart item not found")
	}

	// Get cart item
	cartItem, err := s.cartRepo.GetCartItem(itemID)
	if err != nil {
		logger.Error("Failed to get cart item %d: %v", itemID, err)
		return nil, errors.New("failed to update cart item")
	}
	if cartItem == nil {
		return nil, errors.New("cart item not found")
	}

	// Validate product stock
	product, err := s.productRepo.GetByID(cartItem.ProductID)
	if err != nil {
		logger.Error("Failed to get product %d: %v", cartItem.ProductID, err)
		return nil, errors.New("failed to update cart item")
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if product.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock. Available: %d", product.Quantity)
	}

	// Update quantity
	cartItem.Quantity = req.Quantity
	if err := s.cartRepo.UpdateItem(cartItem); err != nil {
		logger.Error("Failed to update cart item %d: %v", itemID, err)
		return nil, errors.New("failed to update cart item")
	}

	logger.Success("Updated cart item %d quantity to %d for user %d", itemID, req.Quantity, userID)

	// Reload cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		logger.Error("Failed to reload cart: %v", err)
		return nil, errors.New("failed to get updated cart")
	}

	response := cart.ToResponse()
	return &response, nil
}

// ========================================
// REMOVE FROM CART
// ========================================

// RemoveCartItem removes an item from cart
func (s *CartService) RemoveCartItem(userID, itemID uint) error {
	// Check ownership
	isOwner, err := s.cartRepo.IsCartItemOwner(itemID, userID)
	if err != nil {
		logger.Error("Failed to check cart item ownership: %v", err)
		return errors.New("failed to remove cart item")
	}
	if !isOwner {
		return errors.New("cart item not found")
	}

	// Delete item
	if err := s.cartRepo.DeleteItem(itemID); err != nil {
		logger.Error("Failed to delete cart item %d: %v", itemID, err)
		return errors.New("failed to remove cart item")
	}

	logger.Success("Removed cart item %d for user %d", itemID, userID)
	return nil
}

// ========================================
// CLEAR CART
// ========================================

// ClearCart removes all items from user's cart
func (s *CartService) ClearCart(userID uint) error {
	// Get user's cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		logger.Error("Failed to get cart for user %d: %v", userID, err)
		return errors.New("failed to clear cart")
	}
	if cart == nil {
		return errors.New("cart not found")
	}

	// Clear all items
	if err := s.cartRepo.ClearCart(cart.ID); err != nil {
		logger.Error("Failed to clear cart %d: %v", cart.ID, err)
		return errors.New("failed to clear cart")
	}

	logger.Success("Cleared cart for user %d", userID)
	return nil
}

// ========================================
// VALIDATION HELPERS
// ========================================

// ValidateCartForCheckout validates cart before creating order
func (s *CartService) ValidateCartForCheckout(userID uint) (*models.Cart, error) {
	// Get cart
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		logger.Error("Failed to get cart for user %d: %v", userID, err)
		return nil, errors.New("failed to validate cart")
	}
	if cart == nil || len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Validate each item
	for _, item := range cart.Items {
		// Check product exists
		if item.Product == nil {
			return nil, fmt.Errorf("product %d not found in cart", item.ProductID)
		}

		// Check product is active
		if !item.Product.IsActive {
			return nil, fmt.Errorf("product '%s' is no longer available", item.Product.Name)
		}

		// Check stock
		if item.Product.Quantity < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for '%s'. Available: %d, requested: %d",
				item.Product.Name, item.Product.Quantity, item.Quantity)
		}
	}

	return cart, nil
}
