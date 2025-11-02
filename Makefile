.PHONY: help build clean test fmt lint run install dev

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the CLI binary
	@echo "Building goxcel..."
	go build -o goxcel ./cmd
	@echo "Build complete: ./goxcel"

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	GOOS=darwin GOARCH=amd64 go build -o goxcel-darwin-amd64 ./cmd
	GOOS=darwin GOARCH=arm64 go build -o goxcel-darwin-arm64 ./cmd
	GOOS=linux GOARCH=amd64 go build -o goxcel-linux-amd64 ./cmd
	GOOS=windows GOARCH=amd64 go build -o goxcel-windows-amd64.exe ./cmd
	@echo "Multi-platform build complete"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -f goxcel goxcel-* output.xlsx
	@echo "Clean complete"

test: ## Run tests
	@echo "Running tests..."
	go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "Format complete"

fmt-check: ## Check if code is formatted
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files need formatting:"; \
		gofmt -l .; \
		exit 1; \
	else \
		echo "All Go files are properly formatted!"; \
	fi

lint: ## Lint code
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Install with:"; \
		echo "  brew install golangci-lint"; \
		echo "  # or"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

run: build ## Build and run with sample data
	@echo "Running goxcel with sample data..."
	./goxcel generate --template .examples/sample.gxl --data .examples/sample-data.json --output output.xlsx
	@echo "Generated: output.xlsx"

run-dry: build ## Build and run in dry-run mode
	@echo "Running goxcel in dry-run mode..."
	./goxcel generate --template .examples/sample.gxl --data .examples/sample-data.json --dry-run

install: ## Install the binary to $GOPATH/bin
	@echo "Installing goxcel..."
	go install ./cmd
	@echo "Installed to $$(go env GOPATH)/bin/goxcel"

dev: fmt vet test ## Run development checks (format, vet, test)
	@echo "Development checks complete"

check: fmt-check vet lint test ## Run all checks
	@echo "All checks passed!"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated"

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated"

example: ## Run example generation
	@echo "Generating example Excel file..."
	@if [ ! -f goxcel ]; then make build; fi
	./goxcel generate --template .examples/sample.gxl --data .examples/sample-data.json --output output.xlsx
	@echo "Example generated: output.xlsx"

debug: ## Build with debug flags and run
	@echo "Building with debug flags..."
	go build -gcflags="all=-N -l" -o goxcel-debug ./cmd
	@echo "Debug build complete: ./goxcel-debug"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

profile: ## Run with CPU profiling
	@echo "Running with CPU profiling..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	@echo "Profile saved to cpu.prof"
	@echo "View with: go tool pprof cpu.prof"

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t goxcel:latest .
	@echo "Docker image built: goxcel:latest"

.DEFAULT_GOAL := help
