# Go parameters
GO_BIN := go
APP_NAME := thread-pool
BUILD_DIR := ./bin
SRC_DIR := ./cmd/$(APP_NAME)

# Default target
.PHONY: all
all: build

# Build the project
.PHONY: build
build:
	@echo "Building the application..."
	$(GO_BIN) build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

# Run the application
.PHONY: run
run:
	@echo "Running the application..."
	$(GO_BIN) run $(SRC_DIR)

# Test the project
.PHONY: test
test:
	@echo "Running tests..."
	$(GO_BIN) test ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO_BIN) fmt ./...

# Lint the code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build directory..."
	rm -rf $(BUILD_DIR)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO_BIN) mod tidy

# Generate code (e.g., protobufs or mocks)
.PHONY: generate
generate:
	@echo "Generating code..."
	$(GO_BIN) generate ./...

# Help command
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build      Build the application"
	@echo "  make run        Run the application"
	@echo "  make test       Run tests"
	@echo "  make fmt        Format the code"
	@echo "  make lint       Lint the code (requires golangci-lint)"
	@echo "  make clean      Clean build artifacts"
	@echo "  make deps       Install dependencies"
	@echo "  make generate   Run go generate"
	@echo "  make help       Show this help message"
