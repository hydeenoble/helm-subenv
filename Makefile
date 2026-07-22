.PHONY: build test clean install lint fmt vet coverage coverage-report coverage-check help

# Variables
BINARY_NAME=subenv
BUILD_DIR=bin
GO=go
GOFLAGS=-v
VERSION?=$(shell cat plugin.yaml | grep "version" | cut -d '"' -f 2)
LDFLAGS=-ldflags "-X github.com/hydeenoble/helm-env/cmd.version=$(VERSION)"
COVERAGE_THRESHOLD=48

# Default target
all: fmt vet test build

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  build             - Build the binary"
	@echo "  test              - Run tests"
	@echo "  coverage          - Run tests with coverage and generate HTML report"
	@echo "  coverage-report   - Generate coverage HTML report"
	@echo "  coverage-check    - Run tests and verify coverage threshold ($(COVERAGE_THRESHOLD)%)"
	@echo "  clean             - Remove build artifacts"
	@echo "  install           - Install the binary"
	@echo "  lint              - Run golangci-lint"
	@echo "  fmt               - Format code"
	@echo "  vet               - Run go vet"
	@echo "  all               - Run fmt, vet, test, and build"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

## test: Run tests
test:
	@echo "Running tests..."
	$(GO) test $(GOFLAGS) -race ./...

## coverage: Run tests with coverage and generate HTML report
coverage: coverage-check coverage-report
	@echo "✓ Coverage report generated: coverage.html"

## coverage-report: Generate coverage HTML report from existing coverage.out
coverage-report:
	@if [ ! -f coverage.out ]; then \
		echo "Error: coverage.out not found. Run 'make coverage-check' first."; \
		exit 1; \
	fi
	@$(GO) tool cover -html=coverage.out -o coverage.html

## coverage-check: Run tests and verify coverage threshold
coverage-check:
	@echo "Running tests with coverage (threshold: $(COVERAGE_THRESHOLD)%)..."
	@$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@COVERAGE=$$($(GO) tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$COVERAGE%"; \
	if [ "$$(echo "$$COVERAGE < $(COVERAGE_THRESHOLD)" | bc)" -eq 1 ]; then \
		echo "✗ Coverage $$COVERAGE% is below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "✓ Coverage $$COVERAGE% meets threshold"; \
	fi

## clean: Remove build artifacts and coverage files
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -rf coverage/
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
