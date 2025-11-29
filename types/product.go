package types

import "time"

const (
	ProductStatusActive      = "active"
	ProductStatusInactive    = "inactive"
	ProductStatusOutOfStock  = "out_of_stock"
	ProductStatusDiscontinued = "discontinued"
)

type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Image       string    `json:"image" db:"image"`
	Price       float64   `json:"price" db:"price"`
	Quantity    int       `json:"quantity" db:"quantity"` // Note: Not ACID-compliant, see below
	Status      string    `json:"status" db:"status"`     // active, inactive, out_of_stock
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Note about Quantity field:
// This isn't ACID-compliant because concurrent updates can cause race conditions.
// For production, consider:
// 1. Use database transactions with SELECT ... FOR UPDATE
// 2. Use optimistic locking (version field)
// 3. Use a separate inventory service with proper locking

// ProductStore defines the interface for product data operations
type ProductStore interface {
	GetProductByID(id int) (*Product, error)
	GetProductsByID(ids []int) ([]Product, error)
	GetProducts() ([]*Product, error)
	GetProductsByStatus(status string) ([]*Product, error)
	CreateProduct(CreateProductPayload) error
	UpdateProduct(Product) error
	DeleteProduct(id int) error
	UpdateProductQuantity(id int, quantity int) error // Atomic quantity update
}

// CreateProductPayload represents the data needed to create a new product
type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required,min=3,max=200"`
	Description string  `json:"description" validate:"max=2000"`
	Image       string  `json:"image" validate:"omitempty,url"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Quantity    int     `json:"quantity" validate:"required,gte=0"`
	Status      string  `json:"status" validate:"omitempty,oneof=active inactive"`
}

// UpdateProductPayload represents the data that can be updated for a product
type UpdateProductPayload struct {
	Name        *string  `json:"name" validate:"omitempty,min=3,max=200"`
	Description *string  `json:"description" validate:"omitempty,max=2000"`
	Image       *string  `json:"image" validate:"omitempty,url"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
	Quantity    *int     `json:"quantity" validate:"omitempty,gte=0"`
	Status      *string  `json:"status" validate:"omitempty,oneof=active inactive out_of_stock discontinued"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	Products   []*Product `json:"products"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}
