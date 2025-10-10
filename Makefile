# Makefile for Go project

# --- Variables ---
# Application Name
APP_NAME ?= $(shell go list -m)

# Version information can be passed at build time
VERSION ?= 1.0.0
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_DATE ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Go command variables
GO_CMD = go
GO_FLAGS = 
GO_BUILD_FLAGS = -v
GO_TEST_FLAGS = -v -race

# Output directory
OUTPUT_DIR = bin

# LDFLAGS to inject version information
LDFLAGS = -s -w -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT)

# --- Targets ---
.PHONY: all build run test clean tidy help

all: build

# Build the application
# Supports cross-compilation via GOOS and GOARCH environment variables
# Example: GOOS=linux GOARCH=amd64 make build
build:
	@echo "==> Building application..."
	@mkdir -p $(OUTPUT_DIR)
	$(GO_CMD) build $(GO_FLAGS) $(GO_BUILD_FLAGS) -ldflags="$(LDFLAGS)" -o $(OUTPUT_DIR)/$(APP_NAME) .

# Run the application
run:
	@echo "==> Running application..."
	$(GO_CMD) run -ldflags="$(LDFLAGS)" .

# Run all tests
# Example: make test TEST_FLAGS="-v -cover"
test:
	@echo "==> Running tests..."
	$(GO_CMD) test $(GO_TEST_FLAGS) ./...

# Tidy module dependencies
tidy:
	@echo "==> Tidying module dependencies..."
	$(GO_CMD) mod tidy

# Clean up build artifacts
clean:
	@echo "==> Cleaning up..."
	@rm -rf $(OUTPUT_DIR)
	@$(GO_CMD) clean -cache -testcache -modcache

# Help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all        Build the application (default)."
	@echo "  build      Build the application for the current OS/ARCH."
	@echo "             Supports cross-compilation (e.g., GOOS=linux make build)."
	@echo "  run        Build and run the application."
	@echo "  test       Run all unit tests."
	@echo "  tidy       Tidy module dependencies."
	@echo "  clean      Clean up all build artifacts and caches."
	@echo "  help       Show this help message."
