package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/capamir/go-api/services/auth"
	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Public routes
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", h.handleGetProduct).Methods(http.MethodGet)

	// Admin routes (protected with JWT)
	router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProduct, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/products/{productID}", auth.WithJWTAuth(h.handleUpdateProduct, h.userStore)).Methods(http.MethodPut)
	router.HandleFunc("/products/{productID}", auth.WithJWTAuth(h.handleDeleteProduct, h.userStore)).Methods(http.MethodDelete)
}

// handleGetProducts returns a list of all products
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.S.Errorf("Failed to get products: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve products"))
		return
	}

	utils.S.Debugf("Retrieved %d products", len(products))
	utils.WriteJSON(w, http.StatusOK, products)
}

// handleGetProduct returns a single product by ID
func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["productID"]
	if !ok {
		utils.S.Warnf("Missing productID in request")
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	productID, err := strconv.Atoi(str)
	if err != nil {
		utils.S.Warnf("Invalid productID: %s", str)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.S.Warnf("Product %d not found", productID)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	utils.S.Debugf("Retrieved product %d: %s", product.ID, product.Name)
	utils.WriteJSON(w, http.StatusOK, product)
}

// handleCreateProduct creates a new product (admin only)
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.S.Errorf("Failed to get user from context: %v", err)
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
		return
	}

	// Parse payload
	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.S.Warnf("Invalid JSON from user %d: %v", userID, err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.S.Warnf("Validation failed for user %d: %v", userID, validationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrors))
		return
	}

	// Create product
	if err := h.store.CreateProduct(payload); err != nil {
		utils.S.Errorf("Failed to create product by user %d: %v", userID, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to create product"))
		return
	}

	utils.S.Successf("User %d created product: %s", userID, payload.Name)
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Product created successfully",
		"product": payload,
	})
}

// handleUpdateProduct updates an existing product (admin only)
func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	productIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// Check if product exists
	existingProduct, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	// Parse payload
	var payload types.CreateProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrors))
		return
	}

	// Update product
	updatedProduct := types.Product{
		ID:          existingProduct.ID,
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
		CreatedAt:   existingProduct.CreatedAt,
	}

	if err := h.store.UpdateProduct(updatedProduct); err != nil {
		utils.S.Errorf("Failed to update product %d by user %d: %v", productID, userID, err)
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update product"))
		return
	}

	utils.S.Successf("User %d updated product %d", userID, productID)
	utils.WriteJSON(w, http.StatusOK, updatedProduct)
}

// handleDeleteProduct deletes a product (admin only)
func (h *Handler) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("authentication required"))
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	productIDStr, ok := vars["productID"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product ID"))
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// Delete product
	if err := h.store.DeleteProduct(productID); err != nil {
		utils.S.Errorf("Failed to delete product %d by user %d: %v", productID, userID, err)
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	utils.S.Successf("User %d deleted product %d", userID, productID)
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Product deleted successfully",
	})
}
