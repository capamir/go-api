package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// scanRowIntoProduct maps a database row into a Product struct
func scanRowIntoProduct(row interface {
	Scan(dest ...interface{}) error
}) (*types.Product, error) {
	product := new(types.Product)

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}

	return product, nil
}

// GetProductByID retrieves a single product by ID
func (s *Store) GetProductByID(productID int) (*types.Product, error) {
	query := `
		SELECT id, name, description, image, price, quantity, created_at 
		FROM products 
		WHERE id = ?
	`

	row := s.db.QueryRow(query, productID)
	
	product, err := scanRowIntoProduct(row)
	if err == sql.ErrNoRows {
		utils.S.Debugf("Product not found: %d", productID)
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		utils.S.Errorf("Failed to get product %d: %v", productID, err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	utils.S.Debugf("Retrieved product %d: %s", product.ID, product.Name)
	return product, nil
}

// GetProducts retrieves all products
func (s *Store) GetProducts() ([]*types.Product, error) {
	query := `
		SELECT id, name, description, image, price, quantity, created_at 
		FROM products 
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		utils.S.Errorf("Failed to query products: %v", err)
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	products := make([]*types.Product, 0)
	for rows.Next() {
		product := new(types.Product)
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
		)
		if err != nil {
			utils.S.Warnf("Failed to scan product row: %v", err)
			continue
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		utils.S.Errorf("Error iterating product rows: %v", err)
		return nil, fmt.Errorf("error reading products: %w", err)
	}

	utils.S.Debugf("Retrieved %d products", len(products))
	return products, nil
}

// GetProductsByID retrieves multiple products by their IDs
func (s *Store) GetProductsByID(productIDs []int) ([]types.Product, error) {
	if len(productIDs) == 0 {
		return []types.Product{}, nil
	}

	// Build placeholders for IN clause: (?, ?, ?)
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf(`
		SELECT id, name, description, image, price, quantity, created_at 
		FROM products 
		WHERE id IN (?%s)
	`, placeholders)

	// Convert []int to []interface{} for Query variadic args
	args := make([]interface{}, len(productIDs))
	for i, id := range productIDs {
		args[i] = id
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		utils.S.Errorf("Failed to query products by IDs: %v", err)
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	products := make([]types.Product, 0, len(productIDs))
	for rows.Next() {
		product := types.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
		)
		if err != nil {
			utils.S.Warnf("Failed to scan product row: %v", err)
			continue
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		utils.S.Errorf("Error iterating product rows: %v", err)
		return nil, fmt.Errorf("error reading products: %w", err)
	}

	utils.S.Debugf("Retrieved %d products from IDs list", len(products))
	return products, nil
}

// CreateProduct creates a new product
func (s *Store) CreateProduct(payload types.CreateProductPayload) error {
	query := `
		INSERT INTO products (name, description, image, price, quantity, created_at) 
		VALUES (?, ?, ?, ?, ?, NOW())
	`

	res, err := s.db.Exec(
		query,
		payload.Name,
		payload.Description,
		payload.Image,
		payload.Price,
		payload.Quantity,
	)
	if err != nil {
		utils.S.Errorf("Failed to create product '%s': %v", payload.Name, err)
		return fmt.Errorf("failed to create product: %w", err)
	}

	id, err := res.LastInsertId()
	if err == nil {
		utils.S.Successf("Created product %d: %s", id, payload.Name)
	}

	return nil
}

// UpdateProduct updates an existing product
func (s *Store) UpdateProduct(product types.Product) error {
	query := `
		UPDATE products 
		SET name = ?, description = ?, image = ?, price = ?, quantity = ?, updated_at = NOW() 
		WHERE id = ?
	`

	res, err := s.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Image,
		product.Price,
		product.Quantity,
		product.ID,
	)
	if err != nil {
		utils.S.Errorf("Failed to update product %d: %v", product.ID, err)
		return fmt.Errorf("failed to update product: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product %d not found", product.ID)
	}

	utils.S.Successf("Updated product %d: %s", product.ID, product.Name)
	return nil
}

// UpdateProductQuantity atomically updates a product's quantity
func (s *Store) UpdateProductQuantity(productID int, quantity int) error {
	query := `UPDATE products SET quantity = ?, updated_at = NOW() WHERE id = ?`

	res, err := s.db.Exec(query, quantity, productID)
	if err != nil {
		utils.S.Errorf("Failed to update quantity for product %d: %v", productID, err)
		return fmt.Errorf("failed to update product quantity: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product %d not found", productID)
	}

	utils.S.Debugf("Updated product %d quantity to %d", productID, quantity)
	return nil
}

// DeleteProduct soft-deletes or hard-deletes a product
func (s *Store) DeleteProduct(productID int) error {
	// Hard delete
	query := `DELETE FROM products WHERE id = ?`

	// For soft delete, use:
	// query := `UPDATE products SET deleted_at = NOW() WHERE id = ?`

	res, err := s.db.Exec(query, productID)
	if err != nil {
		utils.S.Errorf("Failed to delete product %d: %v", productID, err)
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product %d not found", productID)
	}

	utils.S.Successf("Deleted product %d", productID)
	return nil
}

// GetProductsByStatus retrieves products by their status
func (s *Store) GetProductsByStatus(status string) ([]*types.Product, error) {
	query := `
		SELECT id, name, description, image, price, quantity, created_at 
		FROM products 
		WHERE status = ? 
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, status)
	if err != nil {
		utils.S.Errorf("Failed to query products by status %s: %v", status, err)
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	products := make([]*types.Product, 0)
	for rows.Next() {
		product := new(types.Product)
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.Quantity,
			&product.CreatedAt,
		)
		if err != nil {
			utils.S.Warnf("Failed to scan product row: %v", err)
			continue
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading products: %w", err)
	}

	return products, nil
}
