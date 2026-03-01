.PHONY: all build install clean test test-cover lint fmt vet help build-weather build-ask-cli

# Variables
BINARY_DIR := bin
BINARIES := weather
GO := go
GOFLAGS := -v
GOPATH := $(shell go env GOPATH)

# Default target
all: build

# Build all CLIs
build:
	@echo "Building all CLIs..."
	@mkdir -p $(BINARY_DIR)
	@for binary in $(BINARIES); do \
		echo "Building $$binary..."; \
		$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/$$binary ./cmd/$$binary || exit 1; \
	done
	@echo "Build complete! Binaries are in $(BINARY_DIR)/"

# Build weather CLI only
build-weather:
	@echo "Building weather CLI..."
	@mkdir -p $(BINARY_DIR)
	$(GO) build $(GOFLAGS) -o $(BINARY_DIR)/weather ./cmd/weather

# Install CLIs to system
install: build
	@echo "Installing CLIs to $(GOPATH)/bin..."
	@for binary in $(BINARIES); do \
		echo "Installing $$binary..."; \
		$(GO) install ./cmd/$$binary || exit 1; \
	done
	@echo "Installation complete!"

# Run all tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test -v -cover -coverprofile=coverage.out ./...
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GO) vet ./...

# Run linters
lint: fmt vet
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/"; \
	fi

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  make build          - Build all CLIs (default)"
	@echo "  make build-weather  - Build weather CLI only"
	@echo "  make install        - Install CLIs to system (GOPATH/bin)"
	@echo "  make test           - Run all tests"
	@echo "  make test-cover     - Run tests with coverage report"
	@echo "  make fmt            - Format code with gofmt"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run all linters (fmt, vet, golangci-lint)"
	@echo "  make tidy           - Tidy dependencies"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make help           - Show this help message"
