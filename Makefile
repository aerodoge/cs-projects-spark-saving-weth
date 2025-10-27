# Makefile for cs-projects-spark-saving-weth

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golint
GOVET=$(GOCMD) vet

# Binary name
BINARY_NAME=spark-saving-weth
BINARY_UNIX=$(BINARY_NAME)_unix

# Build directory
BUILD_DIR=build

.PHONY: all build clean test coverage deps fmt lint vet help run

# Default target
all: clean deps fmt lint vet test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v .

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

# Tidy up dependencies
deps-tidy:
	@echo "Tidying up dependencies..."
	$(GOMOD) tidy

# Format Go code
fmt:
	@echo "Formatting Go code..."
	$(GOFMT) -w .

# Check formatting
fmt-check:
	@echo "Checking Go code formatting..."
	@test -z $$($(GOFMT) -l .) || (echo "Code not formatted properly. Run 'make fmt'" && exit 1)

# Lint Go code (requires golint: go install golang.org/x/lint/golint@latest)
lint:
	@echo "Linting Go code..."
	@which golint > /dev/null || (echo "golint not installed. Run: go install golang.org/x/lint/golint@latest" && exit 1)
	golint ./...

# Vet Go code
vet:
	@echo "Vetting Go code..."
	$(GOVET) ./...

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run .

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v .

# Build for multiple platforms
build-all: build build-linux
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe -v .

# Install dependencies for development
dev-deps:
	@echo "Installing development dependencies..."
	$(GOCMD) install golang.org/x/lint/golint@latest
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run golangci-lint (requires golangci-lint)
lint-advanced:
	@echo "Running advanced linting..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run 'make dev-deps'" && exit 1)
	golangci-lint run

# Security scan (requires gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
security:
	@echo "Running security scan..."
	@which gosec > /dev/null || (echo "gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...

# Generate contract bindings (if you have contract ABIs)
generate:
	@echo "Generating contract bindings..."
	@# Add abigen commands here if needed
	@# Example: abigen --abi contract/abis/WETH.abi --pkg contract --type WETH --out contract/weth.go

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, deps, format, lint, vet, test, and build"
	@echo "  build        - Build the binary"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  deps         - Download dependencies"
	@echo "  deps-tidy    - Tidy up dependencies"
	@echo "  fmt          - Format Go code"
	@echo "  fmt-check    - Check if code is formatted"
	@echo "  lint         - Lint Go code"
	@echo "  lint-advanced- Run advanced linting with golangci-lint"
	@echo "  vet          - Vet Go code"
	@echo "  run          - Run the application"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  dev-deps     - Install development dependencies"
	@echo "  security     - Run security scan"
	@echo "  generate     - Generate contract bindings"
	@echo "  docker-build - Build Docker image"
	@echo "  help         - Show this help message"