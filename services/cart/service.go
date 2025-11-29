package cart

import (
	"fmt"

	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
)

// getCartItemsIDs extracts product IDs from cart items and validates quantities
func getCartItemsIDs(items []types.CartCheckoutItem) ([]int, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity %d for product %d", item.Quantity, item.ProductID)
		}
		if item.ProductID <= 0 {
			return nil, fmt.Errorf("invalid product ID: %d", item.ProductID)
		}
		productIds[i] = item.ProductID
	}

	return productIds, nil
}

// checkIfCartIsInStock validates that all cart items are available in sufficient quantity
func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product '%s' only has %d units available, but %d were requested", 
				product.Name, product.Quantity, item.Quantity)
		}

		// Additional validation: check if product is active/available for sale
		// if product.Status != "active" {
		//     return fmt.Errorf("product '%s' is currently unavailable", product.Name)
		// }
	}

	return nil
}

// calculateTotalPrice computes the total price for all items in the cart
func calculateTotalPrice(cartItems []types.CartCheckoutItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal
		
		utils.S.Debugf("Product: %s, Qty: %d, Price: $%.2f, Subtotal: $%.2f", 
			product.Name, item.Quantity, product.Price, itemTotal)
	}

	return total
}

// createOrder processes the cart checkout and creates an order with all items
func (h *Handler) createOrder(products []types.Product, cartItems []types.CartCheckoutItem, userID int) (int, float64, error) {
	// Create a map of products for O(1) lookup
	productsMap := make(map[int]types.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	// Validate: Check if all products are in stock
	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
		return 0, 0, err
	}

	// Calculate total price
	totalPrice := calculateTotalPrice(cartItems, productsMap)
	utils.S.Infof("Cart total for user %d: $%.2f", userID, totalPrice)

	// TODO: In production, wrap this in a database transaction
	// tx, err := h.db.Begin()
	// if err != nil {
	//     return 0, 0, fmt.Errorf("failed to start transaction: %w", err)
	// }
	// defer tx.Rollback() // Will be ignored if tx.Commit() succeeds

	// Step 1: Create the order record
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // TODO: Get from request payload or user profile
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create order: %w", err)
	}

	utils.S.Debugf("Created order %d for user %d", orderID, userID)

	// Step 2: Create order items and update product quantities
	for _, item := range cartItems {
		product := productsMap[item.ProductID]

		// Create order item record
		err := h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price, // Store price at time of purchase
		})
		if err != nil {
			// TODO: In production, this should rollback the transaction
			utils.S.Errorf("Failed to create order item for product %d: %v", item.ProductID, err)
			return 0, 0, fmt.Errorf("failed to create order item: %w", err)
		}

		// Update product quantity (reduce stock)
		product.Quantity -= item.Quantity
		err = h.store.UpdateProduct(product)
		if err != nil {
			// TODO: In production, this should rollback the transaction
			utils.S.Errorf("Failed to update product %d quantity: %v", item.ProductID, err)
			return 0, 0, fmt.Errorf("failed to update product inventory: %w", err)
		}

		utils.S.Debugf("Updated product %d, new quantity: %d", product.ID, product.Quantity)
	}

	// TODO: Commit transaction
	// if err := tx.Commit(); err != nil {
	//     return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	// }

	return orderID, totalPrice, nil
}
