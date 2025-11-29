cat > Makefile << 'EOF'
.PHONY: help run build test clean dev

# Colors
GREEN  := \033[0;32m
YELLOW := \033[0;33m
BLUE   := \033[0;34m
NC     := \033[0m

## help: Show this help message
help:
	@echo "$(BLUE)Available commands:$(NC)"
	@echo "  $(GREEN)make run$(NC)     - Run the application"
	@echo "  $(GREEN)make dev$(NC)     - Run in development mode with auto-reload"
	@echo "  $(GREEN)make build$(NC)   - Build the application"
	@echo "  $(GREEN)make test$(NC)    - Run tests"
	@echo "  $(GREEN)make clean$(NC)   - Clean build artifacts"

## run: Run the application
run:
	@echo "$(BLUE)🚀 Starting application...$(NC)"
	@go run cmd/api/main.go

## dev: Run with air (auto-reload)
dev:
	@echo "$(BLUE)🔥 Starting development mode...$(NC)"
	@air

## build: Build the application
build:
	@echo "$(BLUE)🔨 Building application...$(NC)"
	@go build -o bin/api cmd/api/main.go
	@echo "$(GREEN)✅ Build complete: bin/api$(NC)"

## test: Run tests
test:
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	@go test -v ./...

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)🧹 Cleaning...$(NC)"
	@rm -rf bin/
	@echo "$(GREEN)✅ Clean complete$(NC)"
EOF
