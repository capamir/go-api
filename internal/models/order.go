package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Order status constants
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
)

// Payment status constants
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

// Order represents a customer order
type Order struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderNumber     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	User            *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items           []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	Subtotal        float64        `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Tax             float64        `gorm:"type:decimal(10,2);default:0" json:"tax"`
	ShippingCost    float64        `gorm:"type:decimal(10,2);default:0" json:"shipping_cost"`
	Total           float64        `gorm:"type:decimal(10,2);not null" json:"total"`
	Status          string         `gorm:"type:varchar(20);default:'pending'" json:"status"`
	PaymentStatus   string         `gorm:"type:varchar(20);default:'pending'" json:"payment_status"`
	ShippingAddress datatypes.JSON `gorm:"type:json" json:"shipping_address"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}

// ShippingAddressData represents the shipping address structure
type ShippingAddressData struct {
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip_code"`
	Country     string `json:"country"`
}

// OrderResponse is the DTO for API responses
type OrderResponse struct {
	ID              uint                `json:"id"`
	OrderNumber     string              `json:"order_number"`
	UserID          uint                `json:"user_id"`
	Items           []OrderItemResponse `json:"items"`
	Subtotal        float64             `json:"subtotal"`
	Tax             float64             `json:"tax"`
	ShippingCost    float64             `json:"shipping_cost"`
	Total           float64             `json:"total"`
	Status          string              `json:"status"`
	PaymentStatus   string              `json:"payment_status"`
	ShippingAddress ShippingAddressData `json:"shipping_address"`
	Notes           string              `json:"notes,omitempty"`
	ItemCount       int                 `json:"item_count"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// ToResponse converts Order to OrderResponse
func (o *Order) ToResponse() OrderResponse {
	response := OrderResponse{
		ID:            o.ID,
		OrderNumber:   o.OrderNumber,
		UserID:        o.UserID,
		Items:         []OrderItemResponse{},
		Subtotal:      o.Subtotal,
		Tax:           o.Tax,
		ShippingCost:  o.ShippingCost,
		Total:         o.Total,
		Status:        o.Status,
		PaymentStatus: o.PaymentStatus,
		Notes:         o.Notes,
		ItemCount:     0,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}

	// Parse shipping address
	if o.ShippingAddress != nil {
		var address ShippingAddressData
		if err := o.ShippingAddress.Scan(&address); err == nil {
			response.ShippingAddress = address
		}
	}

	// Convert items
	for _, item := range o.Items {
		itemResponse := item.ToResponse()
		response.Items = append(response.Items, itemResponse)
		response.ItemCount += item.Quantity
	}

	return response
}

// CanBeCancelled checks if order can be cancelled
func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusProcessing
}

// CanBeUpdated checks if order status can be updated
func (o *Order) CanBeUpdated() bool {
	return o.Status != OrderStatusDelivered && o.Status != OrderStatusCancelled
}
