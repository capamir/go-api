# Database migration commands
DB_URL=mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

# Load environment variables from .env
include .env
export

.PHONY: migrate-create migrate-up migrate-down migrate-status migrate-force help

# Create a new migration
# Usage: make migrate-create add-product-table
migrate-create:
	@echo "📝 Creating migration: $(filter-out $@,$(MAKECMDGOALS))"
	@migrate create -ext sql -dir cmd/migrate/migrations -seq $(filter-out $@,$(MAKECMDGOALS))
	@echo "✅ Migration files created!"

# Run all pending migrations
migrate-up:
	@echo "⬆️  Running migrations..."
	@go run cmd/migrate/main.go up

# Rollback the last migration
migrate-down:
	@echo "⬇️  Rolling back migration..."
	@go run cmd/migrate/main.go down

# Show migration status
migrate-status:
	@echo "📊 Migration status:"
	@go run cmd/migrate/main.go status

# Force migration version (use carefully!)
# Usage: make migrate-force 20251125220243
migrate-force:
	@echo "⚠️  Forcing version to: $(filter-out $@,$(MAKECMDGOALS))"
	@go run cmd/migrate/main.go force $(filter-out $@,$(MAKECMDGOALS))

# Show all available commands
help:
	@echo "🚀 Available commands:"
	@echo "  make migrate-create <name>  - Create new migration files"
	@echo "  make migrate-up             - Run all pending migrations"
	@echo "  make migrate-down           - Rollback last migration"
	@echo "  make migrate-status         - Show current migration version"
	@echo "  make migrate-force <version>- Force migration version"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-create add-product-table"
	@echo "  make migrate-up"
	@echo "  make migrate-down"

# Prevent make from treating migration names as targets
%:
	@:
