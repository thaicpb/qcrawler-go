# QCrawler - Advanced Web Crawler Makefile
# ========================================

# Variables
BINARY_NAME=qcrawler
VERSION=1.0.0
BUILD_DIR=build
INSTALL_DIR=/usr/local/bin
GO_VERSION=1.21

# Build information
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[0;37m
NC=\033[0m # No Color

.PHONY: help build install uninstall clean test fmt lint deps check run dev create-build-dir build-all release info

# Default target
all: clean deps build

# Help target
help:
	@echo "$(CYAN)QCrawler v$(VERSION) - Build Commands$(NC)"
	@echo "=================================="
	@echo ""
	@echo "$(GREEN)Basic Commands:$(NC)"
	@echo "  $(YELLOW)make build$(NC)     - Build the binary"
	@echo "  $(YELLOW)make install$(NC)   - Build and install to $(INSTALL_DIR)"
	@echo "  $(YELLOW)make uninstall$(NC) - Remove installed binary"
	@echo "  $(YELLOW)make clean$(NC)     - Clean build artifacts"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  $(YELLOW)make dev$(NC)       - Build and run with example"
	@echo "  $(YELLOW)make test$(NC)      - Run tests"
	@echo "  $(YELLOW)make fmt$(NC)       - Format code"
	@echo "  $(YELLOW)make lint$(NC)      - Run linter"
	@echo "  $(YELLOW)make deps$(NC)      - Download dependencies"
	@echo ""
	@echo "$(GREEN)Cross-platform:$(NC)"
	@echo "  $(YELLOW)make build-all$(NC) - Build for all platforms"
	@echo "  $(YELLOW)make release$(NC)   - Create release archives"
	@echo ""
	@echo "$(GREEN)Usage after install:$(NC)"
	@echo "  $(CYAN)qcrawler -u https://example.com$(NC)"
	@echo "  $(CYAN)qcrawler -u https://example.com -f config.json$(NC)"
	@echo "  $(CYAN)qcrawler -u https://example.com -s \"title=h1,content=.article\"$(NC)"

# Check Go installation
check:
	@echo "$(BLUE)Checking Go installation...$(NC)"
	@go version || (echo "$(RED)Go is not installed. Please install Go $(GO_VERSION)+ from https://golang.org$(NC)" && exit 1)
	@echo "$(GREEN)✓ Go is installed$(NC)"

# Download dependencies
deps: check
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

# Create build directory  
create-build-dir:
	@mkdir -p $(BUILD_DIR)

# Build the binary
build: deps create-build-dir
	@echo "$(BLUE)Building $(BINARY_NAME) v$(VERSION)...$(NC)"
	@CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/qcrawler
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Build for all platforms
build-all: deps create-build-dir
	@echo "$(BLUE)Building for all platforms...$(NC)"
	@echo "$(YELLOW)Building for Linux (amd64)...$(NC)"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/qcrawler
	@echo "$(YELLOW)Building for macOS (amd64)...$(NC)"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/qcrawler
	@echo "$(YELLOW)Building for macOS (arm64)...$(NC)"
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/qcrawler
	@echo "$(YELLOW)Building for Windows (amd64)...$(NC)"
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/qcrawler
	@echo "$(GREEN)✓ Cross-platform build complete$(NC)"

# Install the binary
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME) to $(INSTALL_DIR)...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✓ $(BINARY_NAME) installed successfully$(NC)"
	@echo "$(CYAN)Try: $(BINARY_NAME) --version$(NC)"

# Uninstall the binary
uninstall:
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✓ $(BINARY_NAME) uninstalled$(NC)"

# Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f *.json test*.json
	@echo "$(GREEN)✓ Clean complete$(NC)"

# Format code
fmt:
	@echo "$(BLUE)Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

# Run tests
test: deps
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✓ Tests completed$(NC)"

# Run linter (if golangci-lint is installed)
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

# Development run
dev: build
	@echo "$(BLUE)Running development example...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME) -u https://example.com -v
	@echo "$(GREEN)✓ Development run complete$(NC)"

# Run with custom parameters
run: build
	@echo "$(BLUE)Usage: make run URL=https://example.com [CONFIG=config.json] [OUTPUT=output.json]$(NC)"
	@if [ -z "$(URL)" ]; then \
		echo "$(RED)Error: URL is required. Usage: make run URL=https://example.com$(NC)"; \
		exit 1; \
	fi
	@./$(BUILD_DIR)/$(BINARY_NAME) -u $(URL) $(if $(CONFIG),-f $(CONFIG)) $(if $(OUTPUT),-o $(OUTPUT)) -v

# Create release archives
release: build-all
	@echo "$(BLUE)Creating release archives...$(NC)"
	@mkdir -p $(BUILD_DIR)/release
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip -q release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	@echo "$(GREEN)✓ Release archives created in $(BUILD_DIR)/release/$(NC)"

# Show build info
info:
	@echo "$(CYAN)QCrawler Build Information$(NC)"
	@echo "=========================="
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Install Dir: $(INSTALL_DIR)"
	@echo "Go Version: $(shell go version)"