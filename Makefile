.PHONY: build install test clean release

BINARY_NAME=mirage
VERSION?=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
BUILD_DIR=./bin
INSTALL_DIR?=/usr/local/bin

build:
	@echo "Building ${BINARY_NAME}..."
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mirage
	@echo "✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "✓ Installed to $(INSTALL_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f mirage
	@echo "✓ Clean complete"

release:
	@echo "Creating release with GoReleaser..."
	@goreleaser release --clean

snapshot:
	@echo "Building snapshot release..."
	@goreleaser release --snapshot --clean

dev: build
	@$(BUILD_DIR)/$(BINARY_NAME) start

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies installed"

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"
