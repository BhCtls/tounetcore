# TouNetCore Makefile

.PHONY: build run test clean deps seed-apps seed-admin seed-invite help

# Build the application
build:
	@echo "Building TouNetCore..."
	@go build -o bin/tounetcore cmd/server/main.go
	@echo "✅ Build completed: bin/tounetcore"

# Run the application in development mode
run:
	@echo "Starting TouNetCore server..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "✅ Clean completed"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod tidy
	@echo "✅ Dependencies updated"

# Seed default applications
seed-apps:
	@echo "Seeding default applications..."
	@go run cmd/seed/main.go apps

# Create admin user
seed-admin:
	@echo "Creating admin user..."
	@go run cmd/seed/main.go admin

# Generate invite codes
seed-invite:
	@echo "Generating invite codes..."
	@go run cmd/seed/main.go invite

# Initialize the project (setup database and seed data)
init: deps seed-apps seed-admin seed-invite
	@echo "✅ Project initialized successfully"
	@echo ""
	@echo "🚀 You can now start the server with: make run"
	@echo "👤 Admin credentials: admin / admin123"
	@echo ""

# Development setup
dev-setup: init
	@echo "🎉 Development environment is ready!"

# Production build
prod-build:
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/tounetcore-linux cmd/server/main.go
	@CGO_ENABLED=0 GOOS=darwin go build -ldflags="-w -s" -o bin/tounetcore-darwin cmd/server/main.go
	@CGO_ENABLED=0 GOOS=windows go build -ldflags="-w -s" -o bin/tounetcore-windows.exe cmd/server/main.go
	@echo "✅ Production builds completed"

# Show help
help:
	@echo "TouNetCore - Available commands:"
	@echo ""
	@echo "  build         Build the application"
	@echo "  run           Run the application in development mode"
	@echo "  test          Run tests"
	@echo "  test-coverage Run tests with coverage"
	@echo "  clean         Clean build artifacts"
	@echo "  deps          Download dependencies"
	@echo "  seed-apps     Seed default applications"
	@echo "  seed-admin    Create admin user"
	@echo "  seed-invite   Generate invite codes"
	@echo "  init          Initialize the project (deps + seed data)"
	@echo "  dev-setup     Complete development environment setup"
	@echo "  prod-build    Build for production (multiple platforms)"
	@echo "  help          Show this help message"
	@echo ""

# Default target
.DEFAULT_GOAL := help
