package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/capamir/go-api/cmd/api"
	"github.com/capamir/go-api/configs"
	"github.com/capamir/go-api/db"
	"github.com/capamir/go-api/utils"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize database connection
	database, err := db.NewMySQLStorage(mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		utils.S.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Verify database connection
	if err := initStorage(database); err != nil {
		utils.S.Fatalf("Failed to initialize storage: %v", err)
	}

	// Create API server (FIX: Pass database, not nil!)
	server := api.NewAPIServer("127.0.0.1:"+configs.Envs.Port, database)
	
	// Display server banner
	utils.S.ServerBanner("127.0.0.1:"+configs.Envs.Port, "v1.0.0")

	// Start server in goroutine
	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			utils.S.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with 10 second timeout
	utils.S.Warnf("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		utils.S.Fatalf("Server forced to shutdown: %v", err)
	}

	utils.S.Successf("Server exited properly")
}

func initStorage(db *sql.DB) error {
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	utils.S.DBStatus("connected", "Successfully connected!")
	return nil
}
