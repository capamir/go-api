package repository

import (
	"github.com/capamir/go-api/internal/models"
	"gorm.io/gorm"
)

// CartRepository handles cart database operations
type CartRepository struct {
	db *gorm.DB
}

// NewCartRepository creates a new cart repository
func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

// ========================================
// CART OPERATIONS
// ========================================

// GetOrCreateCart gets user's cart or creates a new one
func (r *CartRepository) GetOrCreateCart(userID uint) (*models.Cart, error) {
	var cart models.Cart

	// Try to find existing cart
	err := r.db.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new cart
			cart = models.Cart{UserID: userID}
			if err := r.db.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}

	return &cart, nil
}

// GetCartByUserID retrieves a cart by user ID with items
func (r *CartRepository) GetCartByUserID(userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product.Category").
		Preload("Items.Product.Tags").
		Where("user_id = ?", userID).
		First(&cart).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

// GetCartByID retrieves a cart by ID with items
func (r *CartRepository) GetCartByID(id uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Items.Product").First(&cart, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

// ClearCart removes all items from a cart
func (r *CartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

// DeleteCart deletes a cart (soft delete)
func (r *CartRepository) DeleteCart(id uint) error {
	return r.db.Delete(&models.Cart{}, id).Error
}

// ========================================
// CART ITEM OPERATIONS
// ========================================

// AddItem adds a new item to cart
func (r *CartRepository) AddItem(item *models.CartItem) error {
	return r.db.Create(item).Error
}

// GetCartItem retrieves a cart item by ID
func (r *CartRepository) GetCartItem(id uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Preload("Product").First(&item, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetCartItemByProductID finds a cart item by cart and product ID
func (r *CartRepository) GetCartItemByProductID(cartID, productID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// UpdateItem updates a cart item
func (r *CartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Save(item).Error
}

// DeleteItem removes an item from cart
func (r *CartRepository) DeleteItem(id uint) error {
	return r.db.Delete(&models.CartItem{}, id).Error
}

// GetCartItems retrieves all items in a cart
func (r *CartRepository) GetCartItems(cartID uint) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.db.Preload("Product").Where("cart_id = ?", cartID).Find(&items).Error
	return items, err
}

// CountCartItems counts items in a cart
func (r *CartRepository) CountCartItems(cartID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.CartItem{}).Where("cart_id = ?", cartID).Count(&count).Error
	return count, err
}

// ========================================
// VALIDATION HELPERS
// ========================================

// IsCartOwner checks if a user owns a cart
func (r *CartRepository) IsCartOwner(cartID, userID uint) (bool, error) {
	var cart models.Cart
	err := r.db.Where("id = ? AND user_id = ?", cartID, userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsCartItemOwner checks if a user owns a cart item
func (r *CartRepository) IsCartItemOwner(itemID, userID uint) (bool, error) {
	var item models.CartItem
	err := r.db.Joins("JOIN carts ON carts.id = cart_items.cart_id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&item).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
