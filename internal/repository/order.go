package repository

import (
	"fmt"
	"time"

	"github.com/capamir/go-api/internal/models"
	"gorm.io/gorm"
)

// OrderRepository handles order database operations
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// ========================================
// ORDER OPERATIONS
// ========================================

// Create creates a new order with items in a transaction
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create order
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID retrieves an order by ID with all relations
func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items.Product").
		Preload("User").
		First(&order, id).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetByOrderNumber retrieves an order by order number
func (r *OrderRepository) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items.Product").
		Preload("User").
		Where("order_number = ?", orderNumber).
		First(&order).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// GetUserOrders retrieves all orders for a user with pagination
func (r *OrderRepository) GetUserOrders(userID uint, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Count total
	if err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := r.db.Preload("Items").
		Where("user_id = ?", userID).
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetAll retrieves all orders with pagination (admin)
func (r *OrderRepository) GetAll(limit, offset int, filters *OrderFilterOptions) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.db.Model(&models.Order{})

	// Apply filters
	if filters != nil {
		if filters.Status != "" {
			query = query.Where("status = ?", filters.Status)
		}
		if filters.PaymentStatus != "" {
			query = query.Where("payment_status = ?", filters.PaymentStatus)
		}
		if filters.UserID != nil {
			query = query.Where("user_id = ?", *filters.UserID)
		}
		if filters.FromDate != nil {
			query = query.Where("created_at >= ?", *filters.FromDate)
		}
		if filters.ToDate != nil {
			query = query.Where("created_at <= ?", *filters.ToDate)
		}
		if filters.MinTotal != nil {
			query = query.Where("total >= ?", *filters.MinTotal)
		}
		if filters.MaxTotal != nil {
			query = query.Where("total <= ?", *filters.MaxTotal)
		}
		if filters.SearchQuery != "" {
			searchPattern := "%" + filters.SearchQuery + "%"
			query = query.Where("order_number LIKE ?", searchPattern)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Preload("Items").
		Preload("User").
		Limit(limit).Offset(offset).
		Order("created_at DESC").
		Find(&orders).Error
	
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// OrderFilterOptions represents filter options for orders
type OrderFilterOptions struct {
	Status        string
	PaymentStatus string
	UserID        *uint
	FromDate      *time.Time
	ToDate        *time.Time
	MinTotal      *float64
	MaxTotal      *float64
	SearchQuery   string
}

// Update updates an order
func (r *OrderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus updates order status
func (r *OrderRepository) UpdateStatus(orderID uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// UpdatePaymentStatus updates payment status
func (r *OrderRepository) UpdatePaymentStatus(orderID uint, paymentStatus string) error {
	return r.db.Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("payment_status", paymentStatus).Error
}

// Delete soft deletes an order
func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

// ========================================
// ORDER NUMBER GENERATION
// ========================================

// GenerateOrderNumber generates a unique order number
func (r *OrderRepository) GenerateOrderNumber() (string, error) {
	// Format: ORD-YYYYMMDD-XXXXX
	// Example: ORD-20251202-00001
	
	date := time.Now().Format("20060102")
	prefix := fmt.Sprintf("ORD-%s-", date)
	
	// Find the latest order for today
	var count int64
	err := r.db.Model(&models.Order{}).
		Where("order_number LIKE ?", prefix+"%").
		Count(&count).Error
	
	if err != nil {
		return "", err
	}
	
	// Increment and format
	orderNumber := fmt.Sprintf("%s%05d", prefix, count+1)
	return orderNumber, nil
}

// OrderNumberExists checks if an order number already exists
func (r *OrderRepository) OrderNumberExists(orderNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Order{}).Where("order_number = ?", orderNumber).Count(&count).Error
	return count > 0, err
}

// ========================================
// VALIDATION HELPERS
// ========================================

// IsOrderOwner checks if a user owns an order
func (r *OrderRepository) IsOrderOwner(orderID, userID uint) (bool, error) {
	var order models.Order
	err := r.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ========================================
// STATISTICS (for admin dashboard)
// ========================================

// GetOrderStats retrieves order statistics
func (r *OrderRepository) GetOrderStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total orders
	var totalOrders int64
	r.db.Model(&models.Order{}).Count(&totalOrders)
	stats["total_orders"] = totalOrders

	// Orders by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	r.db.Model(&models.Order{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusCounts)
	stats["by_status"] = statusCounts

	// Total revenue
	var totalRevenue float64
	r.db.Model(&models.Order{}).
		Where("payment_status = ?", models.PaymentStatusPaid).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalRevenue)
	stats["total_revenue"] = totalRevenue

	// Today's orders
	today := time.Now().Format("2006-01-02")
	var todayOrders int64
	r.db.Model(&models.Order{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayOrders)
	stats["today_orders"] = todayOrders

	return stats, nil
}

// GetRecentOrders retrieves recent orders (for admin dashboard)
func (r *OrderRepository) GetRecentOrders(limit int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("User").
		Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
