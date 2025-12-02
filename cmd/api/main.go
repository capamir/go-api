package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/joho/godotenv"
	"github.com/capamir/go-api/internal/config"
	"github.com/capamir/go-api/internal/database"
	"github.com/capamir/go-api/internal/handler"
	"github.com/capamir/go-api/internal/middleware"
	"github.com/capamir/go-api/internal/repository"
	"github.com/capamir/go-api/internal/service"
	"github.com/capamir/go-api/internal/utils"
	"github.com/capamir/go-api/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// 🆕 Load .env file
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}
	
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run auto-migrations
	if err := database.AutoMigrate(); err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}

	// ========================================
	// Initialize Repositories
	// ========================================
	userRepo := repository.NewUserRepository(database.GetDB())
	categoryRepo := repository.NewCategoryRepository(database.GetDB())
	tagRepo := repository.NewTagRepository(database.GetDB()) 
	productRepo := repository.NewProductRepository(database.GetDB()) 
	cartRepo := repository.NewCartRepository(database.GetDB()) 
	orderRepo := repository.NewOrderRepository(database.GetDB())

	// ========================================
	// Initialize Email Service
	// ========================================
	var emailService utils.EmailService
	if cfg.IsProduction() && cfg.Email.From != "" {
		emailService = utils.NewSMTPEmailService(cfg)
		logger.Info("Using SMTP email service (from: %s)", cfg.Email.From)
	} else {
		emailService = utils.NewConsoleEmailService(cfg.App.URL)
		logger.Info("Using console email service (dev mode)")
	}

	// ========================================
	// Initialize Services
	// ========================================
	authService := service.NewAuthService(userRepo, emailService)
	categoryService := service.NewCategoryService(categoryRepo) 
	tagService := service.NewTagService(tagRepo) 
	productService := service.NewProductService(productRepo, categoryRepo, tagRepo) 
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo)

	// ========================================
	// Initialize Handlers
	// ========================================
	authHandler := handler.NewAuthHandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService) 
	tagHandler := handler.NewTagHandler(tagService)
	productHandler := handler.NewProductHandler(productService) 
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)  

	// Set Gin mode
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// ========================================
	// API v1 Routes
	// ========================================
	v1 := router.Group("/api/v1")
	{
		// ========================================
		// Auth routes (public)
		// ========================================
		auth := v1.Group("/auth")
		{
			// Registration & Login
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Email Verification (public - no auth required)
			auth.GET("/verify", authHandler.VerifyEmail)
			auth.POST("/resend-verification", authHandler.ResendVerification)

			// Protected routes (require authentication)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetProfile)
		}

		// ========================================
		// 🆕 Category routes (public)
		// ========================================
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/active", categoryHandler.GetActiveCategories)
			categories.GET("/root", categoryHandler.GetRootCategories)
			categories.GET("/tree", categoryHandler.GetCategoryTree)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
			categories.GET("/:id/children", categoryHandler.GetCategoryWithChildren)
		}

		// ========================================
		// 🆕 Tag routes (public)
		// ========================================
		tags := v1.Group("/tags")
		{
			tags.GET("", tagHandler.GetAllTags)
			tags.GET("/all", tagHandler.GetAllTagsList)
			tags.GET("/search", tagHandler.SearchTags)
			tags.GET("/:id", tagHandler.GetTagByID)
			tags.GET("/slug/:slug", tagHandler.GetTagBySlug)
		}

		// ========================================
		// 🆕 Product routes (public)
		// ========================================
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts)
			products.GET("/featured", productHandler.GetFeaturedProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/slug/:slug", productHandler.GetProductBySlug)
			products.GET("/:id/related", productHandler.GetRelatedProducts)
		}

		// ========================================
		// 🆕 Cart routes (authenticated)
		// ========================================
		cart := v1.Group("/cart")
		cart.Use(middleware.AuthMiddleware()) // All cart routes require authentication
		{
			cart.GET("", cartHandler.GetUserCart)
			cart.POST("/items", cartHandler.AddToCart)
			cart.PUT("/items/:item_id", cartHandler.UpdateCartItem)
			cart.DELETE("/items/:item_id", cartHandler.RemoveCartItem)
			cart.DELETE("", cartHandler.ClearCart)
		}

		// ========================================
		// 🆕 Order routes (authenticated)
		// ========================================
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware()) // All order routes require authentication
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.PUT("/:id/cancel", orderHandler.CancelOrder)
		}

		// ========================================
		// 🆕 Admin routes (protected)
		// ========================================
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware()) // All admin routes require authentication
		{
			// Category management (admin only)
			adminCategories := admin.Group("/categories")
			{
				adminCategories.POST("", categoryHandler.CreateCategory)
				adminCategories.PUT("/:id", categoryHandler.UpdateCategory)
				adminCategories.DELETE("/:id", categoryHandler.DeleteCategory)
			}
			// 🆕 Tag management (admin only)
			adminTags := admin.Group("/tags")
			{
				adminTags.POST("", tagHandler.CreateTag)
				adminTags.PUT("/:id", tagHandler.UpdateTag)
				adminTags.DELETE("/:id", tagHandler.DeleteTag)
			}
			// 🆕 Product management (admin only)
			adminProducts := admin.Group("/products")
			{
				adminProducts.POST("", productHandler.CreateProduct)
				adminProducts.PUT("/:id", productHandler.UpdateProduct)
				adminProducts.DELETE("/:id", productHandler.DeleteProduct)
				adminProducts.PUT("/:id/stock", productHandler.UpdateStock)
			}
			// 🆕 Order management (admin only)
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", orderHandler.GetAllOrders)
				adminOrders.GET("/stats", orderHandler.GetOrderStats)
				adminOrders.PUT("/:id/status", orderHandler.UpdateOrderStatus)
			}
		}
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🚀 E-Commerce API",
			"version": "1.0.0",
			"docs":    "/api/v1/docs",
		})
	})

	// Create HTTP server
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Print banner
	logger.Log.Banner("E-Commerce API", "v1.0.0", cfg.App.Port)

	// Start server
	go func() {
		logger.Info("Server starting on http://localhost%s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Warn("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Success("Server exited properly")
}
