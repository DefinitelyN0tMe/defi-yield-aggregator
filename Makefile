# DeFi Yield Aggregator - Makefile
# Common development tasks and commands

.PHONY: all build run test clean docker-up docker-down docker-logs lint fmt help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
API_BINARY=api-server
WORKER_BINARY=worker

# Build directories
BUILD_DIR=./bin

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Linker flags
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

all: build

## help: Show this help message
help:
	@echo "DeFi Yield Aggregator - Available targets:"
	@echo ""
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build both API server and worker binaries
build: build-api build-worker

## build-api: Build the API server binary
build-api:
	@echo "Building API server..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(API_BINARY) ./cmd/server

## build-worker: Build the worker binary
build-worker:
	@echo "Building worker..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(WORKER_BINARY) ./cmd/worker

## run-api: Run the API server locally
run-api:
	$(GOCMD) run ./cmd/server

## run-worker: Run the worker locally
run-worker:
	$(GOCMD) run ./cmd/worker

## test: Run all tests
test:
	$(GOTEST) -v -race -cover ./...

## test-coverage: Run tests with coverage report
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linters
lint:
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## fmt: Format code
fmt:
	$(GOCMD) fmt ./...

## vet: Run go vet
vet:
	$(GOCMD) vet ./...

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## clean: Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## docker-up: Start all services with Docker Compose
docker-up:
	docker-compose up -d

## docker-down: Stop all services
docker-down:
	docker-compose down

## docker-logs: View logs from all services
docker-logs:
	docker-compose logs -f

## docker-build: Build Docker images
docker-build:
	docker-compose build

## docker-rebuild: Rebuild and restart services
docker-rebuild: docker-down docker-build docker-up

## docker-clean: Remove all containers, volumes, and images
docker-clean:
	docker-compose down -v --rmi local

## migrate: Run database migrations
migrate:
	@echo "Running migrations..."
	@if command -v psql >/dev/null; then \
		psql -h localhost -U defi -d defi_aggregator -f migrations/001_create_tables.sql; \
	else \
		echo "psql not found. Run migrations via Docker:"; \
		echo "docker-compose exec postgres psql -U defi -d defi_aggregator -f /docker-entrypoint-initdb.d/001_create_tables.sql"; \
	fi

## health: Check service health
health:
	@curl -s http://localhost:3000/api/v1/health | jq .

## api-test: Quick API test
api-test:
	@echo "Testing health endpoint..."
	@curl -s http://localhost:3000/api/v1/health | jq .
	@echo "\nTesting pools endpoint..."
	@curl -s "http://localhost:3000/api/v1/pools?limit=5" | jq '.data | length'
	@echo "\nTesting stats endpoint..."
	@curl -s http://localhost:3000/api/v1/stats | jq '.totalPools'

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## swagger: Generate Swagger documentation
swagger:
	@if command -v swag >/dev/null; then \
		swag init -g cmd/server/main.go -o docs; \
	else \
		echo "swag not installed. Run: make install-tools"; \
	fi

## dev: Start development environment with hot reload
dev:
	@if command -v air >/dev/null; then \
		air -c .air.toml; \
	else \
		echo "air not installed. Run: make install-tools"; \
	fi

## dev-worker: Start worker with hot reload
dev-worker:
	@if command -v air >/dev/null; then \
		air -c .air.worker.toml; \
	else \
		echo "air not installed. Run: make install-tools"; \
	fi

# Default target
.DEFAULT_GOAL := help
