package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	ParentID    *uint          `gorm:"index" json:"parent_id,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	ImageURL    string         `gorm:"type:varchar(500)" json:"image_url"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate hook - called before inserting record
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	// Set default values
	if c.IsActive == false && tx.Statement.Changed("IsActive") == false {
		c.IsActive = true
	}
	return nil
}

// CategoryResponse is the DTO for API responses (simplified, no circular references)
type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ParentID    *uint     `json:"parent_id,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse converts Category to CategoryResponse
func (c *Category) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		ParentID:    c.ParentID,
		ImageURL:    c.ImageURL,
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// CategoryWithChildren includes nested subcategories
type CategoryWithChildren struct {
	CategoryResponse
	Children []CategoryResponse `json:"children,omitempty"`
}

// ToResponseWithChildren converts Category to CategoryWithChildren
func (c *Category) ToResponseWithChildren() CategoryWithChildren {
	resp := CategoryWithChildren{
		CategoryResponse: c.ToResponse(),
		Children:         []CategoryResponse{},
	}

	// Convert children
	for _, child := range c.Children {
		resp.Children = append(resp.Children, child.ToResponse())
	}

	return resp
}
