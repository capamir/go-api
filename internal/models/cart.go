package models

import (
	"time"

	"gorm.io/gorm"
)

// Cart represents a user's shopping cart
type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items     []CartItem     `gorm:"foreignKey:CartID" json:"items,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Cart model
func (Cart) TableName() string {
	return "carts"
}

// CartResponse is the DTO for API responses
type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	Items     []CartItemResponse `json:"items"`
	Subtotal  float64            `json:"subtotal"`
	ItemCount int                `json:"item_count"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ToResponse converts Cart to CartResponse with calculations
func (c *Cart) ToResponse() CartResponse {
	response := CartResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		Items:     []CartItemResponse{},
		Subtotal:  0,
		ItemCount: 0,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	// Convert items and calculate totals
	for _, item := range c.Items {
		itemResponse := item.ToResponse()
		response.Items = append(response.Items, itemResponse)
		response.Subtotal += itemResponse.Subtotal
		response.ItemCount += item.Quantity
	}

	return response
}
