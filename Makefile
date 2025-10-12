# Include environment variables from .env file
include .env
export

.PHONY: setup migrate run help

# Default target
help:
	@echo "Available commands:"
	@echo "  setup   - Start docker compose services (database)"
	@echo "  migrate - Run setup and execute database migrations"
	@echo "  run     - Run the Go application"
	@echo "  help    - Show this help message"

# Setup command: run docker compose up -d --build
setup:
	@echo "Starting docker compose services..."
	docker compose up -d --build

# Migrate command: run setup first, then run migrations
migrate: setup
	@echo "Running database migrations..."
	@echo "Waiting for database to be ready..."
	sleep 5
	migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:5432/$(DB_NAME)?sslmode=disable" up

# Run command: run the Go application
run: migrate
	@echo "Running Go application..."
	go run ./cmd/main.go
