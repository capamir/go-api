package cart

import (
	"fmt"
	"net/http"

	"github.com/capamir/go-api/services/auth"
	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store      types.ProductStore
	orderStore types.OrderStore
	userStore  types.UserStore
}

func NewHandler(
	store types.ProductStore,
	orderStore types.OrderStore,
	userStore types.UserStore,
) *Handler {
	return &Handler{
		store:      store,
		orderStore: orderStore,
		userStore:  userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	// ✅ FIX: Properly handle error from GetUserIDFromContext
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.S.Errorf("Failed to get user ID from context: %v", err)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
		return
	}

	// Parse cart payload
	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.S.Warnf("Invalid JSON payload from user %d: %v", userID, err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate cart structure
	if err := utils.Validate.Struct(cart); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.S.Warnf("Validation failed for user %d: %v", userID, validationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrors))
		return
	}

	// Extract product IDs from cart items
	productIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.S.Warnf("Invalid cart items from user %d: %v", userID, err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Fetch products from database
	products, err := h.store.GetProductsByID(productIds)
	if err != nil {
		utils.S.Errorf("Failed to fetch products for user %d: %v", userID, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to fetch products"))
		return
	}

	// Create order with all validations
	orderID, totalPrice, err := h.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.S.Warnf("Order creation failed for user %d: %v", userID, err)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.S.Successf("Order %d created for user %d - Total: $%.2f", orderID, userID, totalPrice)

	// Return success response
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"success":     true,
		"order_id":    orderID,
		"total_price": totalPrice,
		"message":     "Order created successfully",
	})
}
