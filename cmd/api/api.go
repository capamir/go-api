package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/capamir/go-api/middleware"
	"github.com/capamir/go-api/services/cart"
	"github.com/capamir/go-api/services/order"
	"github.com/capamir/go-api/services/product"
	"github.com/capamir/go-api/services/user"
	"github.com/capamir/go-api/utils"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr   string
	db     *sql.DB
	server *http.Server
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.CORSMiddleware)

	// Root endpoint
	router.HandleFunc("/", s.handleRoot).Methods("GET")

	// Health check endpoint (outside versioned API)
	router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// API v1 routes
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	s.registerServices(subrouter)

	// Configure HTTP server with timeouts
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,  // Time to read request body
		WriteTimeout: 15 * time.Second,  // Time to write response
		IdleTimeout:  60 * time.Second,  // Keep-alive timeout
	}

	utils.S.Infof("Server listening on %s", s.addr)
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *APIServer) Shutdown(ctx context.Context) error {
	if s.server != nil {
		utils.S.Infof("Shutting down server...")
		return s.server.Shutdown(ctx)
	}
	return nil
}

// registerServices initializes all service stores and registers their routes
func (s *APIServer) registerServices(subrouter *mux.Router) {
	// Initialize all stores
	userStore := user.NewStore(s.db)
	productStore := product.NewStore(s.db)
	orderStore := order.NewStore(s.db)

	// Register service routes
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	utils.S.Successf("All routes registered successfully")
}

// handleHealth checks database connectivity and returns server status
func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if err := s.db.Ping(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	// Return healthy status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "E-commerce API",
		"version": "v1",
		"endpoints": map[string]string{
			"health": "/health",
			"api":    "/api/v1",
		},
	})
}