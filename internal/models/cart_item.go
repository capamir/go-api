package models

import (
	"time"

	"gorm.io/gorm"
)

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CartID    uint           `gorm:"index;not null" json:"cart_id"`
	Cart      *Cart          `gorm:"foreignKey:CartID" json:"cart,omitempty"`
	ProductID uint           `gorm:"index;not null" json:"product_id"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int            `gorm:"not null;default:1" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"` // Price at time of adding
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for CartItem model
func (CartItem) TableName() string {
	return "cart_items"
}

// BeforeCreate hook - sets price from product if not already set
func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	// Set default quantity
	if ci.Quantity <= 0 {
		ci.Quantity = 1
	}
	return nil
}

// CartItemResponse is the DTO for API responses
type CartItemResponse struct {
	ID        uint            `json:"id"`
	ProductID uint            `json:"product_id"`
	Product   *ProductResponse `json:"product,omitempty"`
	Quantity  int             `json:"quantity"`
	Price     float64         `json:"price"`
	Subtotal  float64         `json:"subtotal"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ToResponse converts CartItem to CartItemResponse
func (ci *CartItem) ToResponse() CartItemResponse {
	response := CartItemResponse{
		ID:        ci.ID,
		ProductID: ci.ProductID,
		Quantity:  ci.Quantity,
		Price:     ci.Price,
		Subtotal:  ci.Price * float64(ci.Quantity),
		CreatedAt: ci.CreatedAt,
		UpdatedAt: ci.UpdatedAt,
	}

	// Include product info if loaded
	if ci.Product != nil {
		productResp := ci.Product.ToResponse()
		response.Product = &productResp
	}

	return response
}
