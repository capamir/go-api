package types

import "time"

type CartCheckoutItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items []CartCheckoutItem `json:"items" validate:"required"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}
