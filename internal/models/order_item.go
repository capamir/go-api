package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OrderID     uint           `gorm:"index;not null" json:"order_id"`
	Order       *Order         `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	ProductID   uint           `gorm:"index;not null" json:"product_id"`
	Product     *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName string         `gorm:"type:varchar(255);not null" json:"product_name"` // Snapshot
	ProductSKU  string         `gorm:"type:varchar(100)" json:"product_sku"`           // Snapshot
	Quantity    int            `gorm:"not null" json:"quantity"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Price at time of order
	Subtotal    float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}

// BeforeCreate hook - calculates subtotal
func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	oi.Subtotal = oi.Price * float64(oi.Quantity)
	return nil
}

// OrderItemResponse is the DTO for API responses
type OrderItemResponse struct {
	ID          uint            `json:"id"`
	ProductID   uint            `json:"product_id"`
	ProductName string          `json:"product_name"`
	ProductSKU  string          `json:"product_sku,omitempty"`
	Product     *ProductResponse `json:"product,omitempty"`
	Quantity    int             `json:"quantity"`
	Price       float64         `json:"price"`
	Subtotal    float64         `json:"subtotal"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ToResponse converts OrderItem to OrderItemResponse
func (oi *OrderItem) ToResponse() OrderItemResponse {
	response := OrderItemResponse{
		ID:          oi.ID,
		ProductID:   oi.ProductID,
		ProductName: oi.ProductName,
		ProductSKU:  oi.ProductSKU,
		Quantity:    oi.Quantity,
		Price:       oi.Price,
		Subtotal:    oi.Subtotal,
		CreatedAt:   oi.CreatedAt,
	}

	// Include product info if loaded (for enriched responses)
	if oi.Product != nil {
		productResp := oi.Product.ToResponse()
		response.Product = &productResp
	}

	return response
}
