# Check status
go run cmd/migrate/main.go status

# Run all pending migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down

# Fix dirty state (if migration failed)
go run cmd/migrate/main.go force 20251125220243
