package order

import (
	"database/sql"
	"fmt"

	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateOrder creates a new order record and returns the order ID
func (s *Store) CreateOrder(order types.Order) (int, error) {
	// Note: Column names should match your actual database schema
	// Assuming snake_case: user_id, total, status, address, created_at
	query := `
		INSERT INTO orders (user_id, total, status, address, created_at) 
		VALUES (?, ?, ?, ?, NOW())
	`

	res, err := s.db.Exec(query, order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		utils.S.Errorf("Failed to insert order for user %d: %v", order.UserID, err)
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		utils.S.Errorf("Failed to get last insert ID for order: %v", err)
		return 0, fmt.Errorf("failed to retrieve order ID: %w", err)
	}

	utils.S.Debugf("Created order %d for user %d with total $%.2f", id, order.UserID, order.Total)
	return int(id), nil
}

// CreateOrderItem creates an order item record linked to an order
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	// Note: Assuming snake_case column names: order_id, product_id, quantity, price
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, price) 
		VALUES (?, ?, ?, ?)
	`

	_, err := s.db.Exec(query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	if err != nil {
		utils.S.Errorf("Failed to create order item (order: %d, product: %d): %v", 
			orderItem.OrderID, orderItem.ProductID, err)
		return fmt.Errorf("failed to create order item: %w", err)
	}

	utils.S.Debugf("Added product %d (qty: %d) to order %d", 
		orderItem.ProductID, orderItem.Quantity, orderItem.OrderID)
	return nil
}

// GetOrderByID retrieves an order by its ID
func (s *Store) GetOrderByID(orderID int) (*types.Order, error) {
	query := `
		SELECT id, user_id, total, status, address, created_at 
		FROM orders 
		WHERE id = ?
	`

	var order types.Order
	err := s.db.QueryRow(query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.Status,
		&order.Address,
		&order.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order %d not found", orderID)
	}
	if err != nil {
		utils.S.Errorf("Failed to fetch order %d: %v", orderID, err)
		return nil, fmt.Errorf("failed to retrieve order: %w", err)
	}

	return &order, nil
}

// GetOrdersByUserID retrieves all orders for a specific user
func (s *Store) GetOrdersByUserID(userID int) ([]types.Order, error) {
	query := `
		SELECT id, user_id, total, status, address, created_at 
		FROM orders 
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		utils.S.Errorf("Failed to fetch orders for user %d: %v", userID, err)
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var order types.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Total,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
		)
		if err != nil {
			utils.S.Errorf("Failed to scan order row: %v", err)
			continue
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// GetOrderItems retrieves all items for a specific order
func (s *Store) GetOrderItems(orderID int) ([]types.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, price 
		FROM order_items 
		WHERE order_id = ?
	`

	rows, err := s.db.Query(query, orderID)
	if err != nil {
		utils.S.Errorf("Failed to fetch order items for order %d: %v", orderID, err)
		return nil, fmt.Errorf("failed to retrieve order items: %w", err)
	}
	defer rows.Close()

	var items []types.OrderItem
	for rows.Next() {
		var item types.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			utils.S.Errorf("Failed to scan order item row: %v", err)
			continue
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order items: %w", err)
	}

	return items, nil
}

// UpdateOrderStatus updates the status of an order
func (s *Store) UpdateOrderStatus(orderID int, status string) error {
	query := `UPDATE orders SET status = ? WHERE id = ?`

	res, err := s.db.Exec(query, status, orderID)
	if err != nil {
		utils.S.Errorf("Failed to update order %d status to %s: %v", orderID, status, err)
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("order %d not found", orderID)
	}

	utils.S.Infof("Updated order %d status to '%s'", orderID, status)
	return nil
}
