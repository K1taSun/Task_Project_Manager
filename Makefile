# Task Project Manager Makefile

# Variables
BINARY_NAME=task-manager
BUILD_DIR=build
DIST_DIR=dist
TEST_DIR=tests
COVERAGE_DIR=coverage

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Task Project Manager - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

.PHONY: run
run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run ./cmd/server

.PHONY: dev
dev: ## Run in development mode with hot reload
	@echo "Running in development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -f $(BINARY_NAME)

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	$(GOTEST) ./tests/unit/...

.PHONY: test-integration
test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	$(GOTEST) ./tests/integration/...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCMD) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "Running tests with race detection..."
	$(GOTEST) -race ./...

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: fmt
fmt: ## Format code with go fmt
	@echo "Formatting code..."
	$(GOFMT) ./...

.PHONY: fmt-check
fmt-check: ## Check if code is formatted correctly
	@echo "Checking code formatting..."
	@if [ "$$(go fmt ./... | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted correctly. Run 'make fmt' to fix."; \
		exit 1; \
	else \
		echo "Code is formatted correctly."; \
	fi

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

.PHONY: install
install: ## Install the application
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/server
	@echo "Installation complete. Run './$(BINARY_NAME)' to start the application."

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):latest

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi $(BINARY_NAME):latest 2>/dev/null || true

.PHONY: release
release: ## Build release binaries
	@echo "Building release binaries..."
	@mkdir -p $(DIST_DIR)
	
	# Build for different platforms
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/server
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/server
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/server
	
	@echo "Release binaries created in $(DIST_DIR)/"

.PHONY: check
check: ## Run all checks (fmt, vet, lint, test)
	@echo "Running all checks..."
	@make fmt-check
	@make vet
	@make lint
	@make test

.PHONY: pre-commit
pre-commit: ## Run pre-commit checks
	@echo "Running pre-commit checks..."
	@make fmt
	@make vet
	@make test

.PHONY: ci
ci: ## Run CI pipeline
	@echo "Running CI pipeline..."
	@make deps
	@make fmt-check
	@make vet
	@make lint
	@make test-coverage
	@make test-race
	@make build

.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v godoc > /dev/null; then \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Installing..."; \
		go install golang.org/x/tools/cmd/godoc@latest; \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	fi

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o docs/swagger; \
	else \
		echo "swag not found. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/server/main.go -o docs/swagger; \
	fi

.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

.PHONY: all
all: clean deps test build ## Clean, install deps, test, and build
