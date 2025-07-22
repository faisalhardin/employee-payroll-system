# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Run project initialization
init:
	cp .env.example .env
	cp files/env/envconfig.yaml.example files/env/envconfig.yaml 

# Setup environment and run docker
docker-setup:
	@echo "Setting up environment..."
	@if [ ! -f .env ]; then \
		echo "Copying .env file..."; \
		cp .env.example .env; \
	fi
	@if [ ! -f files/env/envconfig.yaml ]; then \
		echo "Copying envconfig.yaml file..."; \
		cp files/env/envconfig.yaml.example files/env/envconfig.yaml; \
		echo "Updating database host in config..."; \
		sed -i.bak 's/localhost/psql_bp/g' files/env/envconfig.yaml && rm -f files/env/envconfig.yaml.bak; \
	fi
	@echo "Environment setup complete."

# Create and run containers
docker-run: docker-setup
	@echo "Starting containers..."
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Build containers without running
docker-build: docker-setup
	@echo "Building containers..."
	@if docker compose build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose build; \
	fi

# Run containers in detached mode
docker-run-detached: docker-setup
	@echo "Starting containers in detached mode..."
	@if docker compose up -d --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v -cover ./...
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

# Show help
help:
	@echo "Available commands:"
	@echo "  make docker-setup      - Copy and configure environment files"
	@echo "  make docker-build      - Build Docker containers without running them"
	@echo "  make docker-run        - Set up environment and run containers (interactive mode)"
	@echo "  make docker-run-detached - Set up environment and run containers (detached mode)"
	@echo "  make docker-down       - Stop and remove containers"
	@echo "  make clean            - Remove built binaries"
	@echo "  make test             - Run tests"
	@echo "  make itest            - Run integration tests"
	@echo "  make watch            - Run with live reload (requires air)"

.PHONY: all build run test clean watch docker-run docker-down itest docker-setup docker-build docker-run-detached help

.DEFAULT_GOAL := help
