package types

import "time"

// Order status constants
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusShipped   = "shipped"
	OrderStatusDelivered = "delivered"
	OrderStatusCancelled = "cancelled"
)

type CartCheckoutItem struct {
	ProductID int `json:"product_id" validate:"required,gt=0"`
	Quantity  int `json:"quantity" validate:"required,gt=0,lte=100"` // Max 100 per item
}

type CartCheckoutPayload struct {
	Items   []CartCheckoutItem `json:"items" validate:"required,min=1,dive"`
	Address string             `json:"address" validate:"required,min=10,max=500"`
}

type Order struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Total     float64   `json:"total" db:"total"`
	Status    string    `json:"status" db:"status"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type OrderItem struct {
	ID        int       `json:"id" db:"id"`
	OrderID   int       `json:"order_id" db:"order_id"`
	ProductID int       `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	Price     float64   `json:"price" db:"price"` 
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type OrderWithItems struct {
	Order
	Items []OrderItem `json:"items"`
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
	GetOrderByID(orderID int) (*Order, error)
	GetOrdersByUserID(userID int) ([]Order, error)
	GetOrderItems(orderID int) ([]OrderItem, error)
	UpdateOrderStatus(orderID int, status string) error
}
