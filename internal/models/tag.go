package models

import (
	"time"

	"gorm.io/gorm"
)

// Tag represents a product tag/label
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Slug      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Tag model
func (Tag) TableName() string {
	return "tags"
}

// TagResponse is the DTO for API responses
type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Tag to TagResponse
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
