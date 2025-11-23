package db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	// First, connect without specifying a database
	cfgWithoutDB := cfg
	cfgWithoutDB.DBName = ""
	
	db, err := sql.Open("mysql", cfgWithoutDB.FormatDSN())
	if err != nil {
		return nil, err
	}

	// Test the connection without database
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Check if database exists, if not create it
	if err := createDatabaseIfNotExists(db, cfg.DBName); err != nil {
		return nil, err
	}

	// Close the connection without database
	db.Close()

	// Now connect to the specific database
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// createDatabaseIfNotExists creates the database if it doesn't exist
func createDatabaseIfNotExists(db *sql.DB, dbName string) error {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "`")
	if err != nil {
		return fmt.Errorf("failed to create database %s: %v", dbName, err)
	}
	
	_, err = db.Exec("USE `" + dbName + "`")
	if err != nil {
		return fmt.Errorf("failed to select database %s: %v", dbName, err)
	}
	
	return nil
}