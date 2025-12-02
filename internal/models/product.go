package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Product represents a product in the store
type Product struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description  string         `gorm:"type:text" json:"description"`
	Price        float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	ComparePrice float64        `gorm:"type:decimal(10,2)" json:"compare_price"` // Original price (for discounts)
	SKU          string         `gorm:"type:varchar(100);uniqueIndex" json:"sku"`
	Barcode      string         `gorm:"type:varchar(100)" json:"barcode"`
	Quantity     int            `gorm:"default:0" json:"quantity"`
	Images       datatypes.JSON `gorm:"type:json" json:"images"` // Array of image URLs
	CategoryID   uint           `gorm:"index" json:"category_id"`
	Category     *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags         []Tag          `gorm:"many2many:product_tags;" json:"tags,omitempty"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	IsFeatured   bool           `gorm:"default:false" json:"is_featured"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// BeforeCreate hook - called before inserting record
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	// Set default values
	if !p.IsActive && !tx.Statement.Changed("IsActive") {

		p.IsActive = true
	}
	return nil
}

// ProductResponse is the DTO for API responses
type ProductResponse struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	Slug         string           `json:"slug"`
	Description  string           `json:"description,omitempty"`
	Price        float64          `json:"price"`
	ComparePrice float64          `json:"compare_price,omitempty"`
	Discount     float64          `json:"discount,omitempty"` // Calculated discount percentage
	SKU          string           `json:"sku,omitempty"`
	Barcode      string           `json:"barcode,omitempty"`
	Quantity     int              `json:"quantity"`
	InStock      bool             `json:"in_stock"`
	Images       []string         `json:"images"`
	Category     *CategoryResponse `json:"category,omitempty"`
	Tags         []TagResponse    `json:"tags,omitempty"`
	IsActive     bool             `json:"is_active"`
	IsFeatured   bool             `json:"is_featured"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// ToResponse converts Product to ProductResponse
func (p *Product) ToResponse() ProductResponse {
	response := ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Description:  p.Description,
		Price:        p.Price,
		ComparePrice: p.ComparePrice,
		SKU:          p.SKU,
		Barcode:      p.Barcode,
		Quantity:     p.Quantity,
		InStock:      p.Quantity > 0,
		Images:       []string{},
		IsActive:     p.IsActive,
		IsFeatured:   p.IsFeatured,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	// Calculate discount percentage
	if p.ComparePrice > 0 && p.Price < p.ComparePrice {
		response.Discount = ((p.ComparePrice - p.Price) / p.ComparePrice) * 100
	}

	// Parse images from JSON
	if p.Images != nil {
		var images []string
		if err := p.Images.Scan(&images); err == nil {
			response.Images = images
		}
	}

	// Include category if loaded
	if p.Category != nil {
		cat := p.Category.ToResponse()
		response.Category = &cat
	}

	// Include tags if loaded
	if len(p.Tags) > 0 {
		response.Tags = make([]TagResponse, len(p.Tags))
		for i, tag := range p.Tags {
			response.Tags[i] = tag.ToResponse()
		}
	}

	return response
}
