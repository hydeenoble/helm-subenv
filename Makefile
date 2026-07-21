.PHONY: build test clean install lint fmt vet coverage help

# Variables
BINARY_NAME=subenv
BUILD_DIR=bin
GO=go
GOFLAGS=-v
VERSION?=$(shell cat plugin.yaml | grep "version" | cut -d '"' -f 2)
LDFLAGS=-ldflags "-X github.com/hydeenoble/helm-env/cmd.version=$(VERSION)"

# Default target
all: fmt vet test build

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  test       - Run tests"
	@echo "  coverage   - Run tests with coverage"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install the binary"
	@echo "  lint       - Run golangci-lint"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  all        - Run fmt, vet, test, and build"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) ./...

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -rf releases

## install: Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) .

## lint: Run golangci-lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install it from https://golangci-lint.run/usage/install/"; \
	fi

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying go modules..."
	$(GO) mod tidy

## mod-download: Download go modules
mod-download:
	@echo "Downloading go modules..."
	$(GO) mod download
