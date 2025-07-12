# Makefile for Carrion Language Server Protocol (LSP)

# Variables
BINARY_NAME=carrion-lsp
BUILD_DIR=build
GO=go
GOFMT=gofmt
GOLINT=golint
GOVET=go vet

# Default target
.PHONY: all
all: clean fmt vet lint build

# Build the LSP server
.PHONY: build
build:
	@echo "Building Carrion LSP server..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-windows build-macos

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: build-macos
build-macos:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

# Run the LSP server in stdio mode
.PHONY: run
run: build
	@echo "Running Carrion LSP server in stdio mode..."
	$(BUILD_DIR)/$(BINARY_NAME) --stdio

# Run the LSP server in TCP mode
.PHONY: run-tcp
run-tcp: build
	@echo "Running Carrion LSP server in TCP mode on port 7777..."
	$(BUILD_DIR)/$(BINARY_NAME)

# Install to /usr/local/bin
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Test the code
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./...

# Test with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GO) test -race ./...

# Generate test coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Vet the code
.PHONY: vet
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Lint the code
.PHONY: lint
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, install it from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Update dependencies
.PHONY: update
update:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy

# Check for module issues
.PHONY: check
check: fmt vet lint test

# Development mode - build and run with debug logging
.PHONY: dev
dev: build
	@echo "Running in development mode with debug logging..."
	$(BUILD_DIR)/$(BINARY_NAME) --stdio --log=/tmp/carrion-lsp-debug.log

# Create release packages
.PHONY: release
release: clean build-all
	@echo "Creating release packages..."
	@mkdir -p $(BUILD_DIR)/release
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	cd $(BUILD_DIR) && zip release/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@echo "Release packages created in $(BUILD_DIR)/release/"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, format, vet, lint, and build"
	@echo "  build        - Build the LSP server"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  build-macos  - Build for macOS"
	@echo "  run          - Run in stdio mode"
	@echo "  run-tcp      - Run in TCP mode"
	@echo "  install      - Install to /usr/local/bin"
	@echo "  test         - Run tests"
	@echo "  test-race    - Run tests with race detection"
	@echo "  test-coverage- Generate test coverage report"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Lint code"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download dependencies"
	@echo "  update       - Update dependencies"
	@echo "  check        - Run all checks (fmt, vet, lint, test)"
	@echo "  dev          - Run in development mode"
	@echo "  release      - Create release packages"
	@echo "  help         - Show this help"