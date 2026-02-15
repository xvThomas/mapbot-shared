# Makefile for mapbot-shared library

.PHONY: help build test lint test-coverage clean deps

# Colors for output
COLOR_RESET = \033[0m
COLOR_BOLD = \033[1m
COLOR_GREEN = \033[32m
COLOR_YELLOW = \033[33m
COLOR_BLUE = \033[34m

help: ## Show help
	@echo "$(COLOR_BOLD)mapbot-shared - Shared Go utilities$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_GREEN)Available commands:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_BLUE)%-24s$(COLOR_RESET) %s\n", $$1, $$2}'

build: ## Build all packages
	@echo "$(COLOR_YELLOW)Building packages...$(COLOR_RESET)"
	@go build ./...
	@echo "$(COLOR_GREEN)✓ Build successful$(COLOR_RESET)"

test: ## Run tests
	@echo "$(COLOR_YELLOW)Running tests...$(COLOR_RESET)"
	@go test -v ./...

lint: ## Run golangci-lint
	@echo "$(COLOR_YELLOW)Running linter...$(COLOR_RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "$(COLOR_YELLOW)golangci-lint not installed. Install with:$(COLOR_RESET)"; \
		echo "  brew install golangci-lint  # macOS"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

test-coverage: ## Run tests with coverage report
	@echo "$(COLOR_YELLOW)Running tests with coverage...$(COLOR_RESET)"
	@go test -v -coverprofile=coverage.out ./...
	@echo ""
	@echo "=== Coverage summary ==="
	@go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "To view detailed HTML coverage report, run:"
	@echo "  make coverage"
	@echo "Coverage report saved to coverage.out"

coverage: ## Show detailed HTML coverage report
	@echo "$(COLOR_YELLOW)Showing detailed HTML coverage report...$(COLOR_RESET)"
	@go tool cover -html=coverage.out

deps: ## Install dependencies
	@echo "$(COLOR_YELLOW)Installing Go dependencies...$(COLOR_RESET)"
	@go mod download
	@go mod tidy
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

clean: ## Clean generated files
	@echo "$(COLOR_YELLOW)Cleaning...$(COLOR_RESET)"
	@rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✓ Cleaned$(COLOR_RESET)"

all: deps lint test build ## Run all checks and build

.DEFAULT_GOAL := help
