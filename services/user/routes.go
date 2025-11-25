package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/capamir/go-api/services/auth"
	"github.com/capamir/go-api/types"
	"github.com/capamir/go-api/utils"
	"github.com/gorilla/mux"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Route registration logic goes here
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Handle login logic here
}

func (h *Handler) handleRegister(res http.ResponseWriter, req *http.Request) {
	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(req, &payload); err != nil {
		utils.WriteError(res, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(res, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}
	
	// check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		// User exists, return error
		utils.WriteError(res, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}
	
	// User doesn't exist, create new user
	// TODO: Hash the password before storing
	// hash password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(res, http.StatusInternalServerError, err)
		return
	}
	user := types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword, // store hashed password
	}
	
	if err := h.store.CreateUser(user); err != nil {
		utils.WriteError(res, http.StatusInternalServerError, err)
		return
	}
	
	utils.WriteJSON(res, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	// Convert userID to int
	id, err := strconv.Atoi(userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID: %v", err))
		return
	}

	user, err := h.store.GetUserByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)
}
