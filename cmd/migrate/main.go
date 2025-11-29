package main

import (
	"fmt"
	"os"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/capamir/go-api/configs"
	"github.com/capamir/go-api/db"
	"github.com/capamir/go-api/utils"
)

func main() {
	utils.S.Infof("🔄 Database Migration Tool")
	utils.S.Infof("Database: %s", configs.Envs.DBName)

	// Connect to database
	cfg := mysqlDriver.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	database, err := db.NewMySQLStorage(cfg)
	if err != nil {
		utils.S.Fatalf("Failed to connect to database: %v", err)
	}

	// Create migration driver
	driver, err := mysqlMigrate.WithInstance(database, &mysqlMigrate.Config{})
	if err != nil {
		utils.S.Fatalf("Failed to create migration driver: %v", err)
	}

	// Create migrator
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		utils.S.Fatalf("Failed to create migrator: %v", err)
	}

	// Get current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		utils.S.Errorf("Failed to get version: %v", err)
	}

	if dirty {
		utils.S.Warnf("⚠️  Database is in dirty state at version %d", version)
	} else if err == migrate.ErrNilVersion {
		utils.S.Infof("📊 No migrations applied yet")
	} else {
		utils.S.Infof("📊 Current version: %d", version)
	}

	// Parse command
	if len(os.Args) < 2 {
		utils.S.Errorf("Usage: go run cmd/migrate/main.go [up|down|status|force VERSION]")
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		utils.S.Infof("⬆️  Running migrations up...")
		if err := m.Up(); err != nil {
			if err == migrate.ErrNoChange {
				utils.S.Successf("✅ Database is already up to date")
			} else {
				utils.S.Fatalf("Migration up failed: %v", err)
			}
		} else {
			newVersion, _, _ := m.Version()
			utils.S.Successf("✅ Migrations applied successfully! New version: %d", newVersion)
		}

	case "down":
		utils.S.Warnf("⬇️  Running migrations down (rolling back)...")
		if err := m.Down(); err != nil {
			if err == migrate.ErrNoChange {
				utils.S.Infof("No migrations to roll back")
			} else {
				utils.S.Fatalf("Migration down failed: %v", err)
			}
		} else {
			utils.S.Successf("✅ Rollback completed successfully!")
		}

	case "status":
		if dirty {
			utils.S.ErrorBox("Dirty State", fmt.Sprintf("Database is in dirty state at version %d", version))
		} else if err == migrate.ErrNilVersion {
			utils.S.InfoBox("Status", "No migrations applied")
		} else {
			utils.S.SuccessBox("Status", fmt.Sprintf("Current version: %d", version))
		}

	case "force":
		if len(os.Args) < 3 {
			utils.S.Errorf("Usage: go run cmd/migrate/main.go force VERSION")
			os.Exit(1)
		}
		var forceVersion int
		fmt.Sscanf(os.Args[2], "%d", &forceVersion)
		
		utils.S.Warnf("⚠️  Forcing version to %d...", forceVersion)
		if err := m.Force(forceVersion); err != nil {
			utils.S.Fatalf("Force failed: %v", err)
		}
		utils.S.Successf("✅ Version forced to %d", forceVersion)

	default:
		utils.S.Errorf("Unknown command: %s", cmd)
		utils.S.Infof("Available commands: up, down, status, force")
		os.Exit(1)
	}
}
