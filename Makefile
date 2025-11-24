# -----------------------------------------------------------------------------
# Makefile for go-fluent-validator
# -----------------------------------------------------------------------------
# Bu dosya, yaygın geliştirme görevleri için kısayollar sağlar.
#
# Kullanım:
#   make test          - Run all tests
#   make test-verbose  - Run tests with verbose output
#   make coverage      - Generate coverage report
#   make lint          - Run linters
#   make fmt           - Format code
#   make bench         - Run benchmarks
#   make examples      - Run examples
#   make ci            - Run all CI checks locally
#   make clean         - Clean build artifacts
#
# Metadata:
# @author   Ahmet ALTUN
# @github   github.com/biyonik
# @linkedin linkedin.com/in/biyonik
# @email    ahmet.altun60@gmail.com
# -----------------------------------------------------------------------------

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -timeout 5m ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running verbose tests..."
	@go test -v -race -timeout 5m -count=1 ./...

.PHONY: test-short
test-short: ## Run short tests only
	@echo "Running short tests..."
	@go test -short -race ./...

.PHONY: coverage
coverage: ## Generate test coverage report
	@echo "Generating coverage report..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

.PHONY: coverage-ci
coverage-ci: ## Generate coverage for CI (no HTML)
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

.PHONY: bench
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem -run=^$ ./...

.PHONY: bench-compare
bench-compare: ## Run benchmarks and save results
	@echo "Running benchmarks and saving results..."
	@go test -bench=. -benchmem -run=^$ ./... | tee benchmark.txt

.PHONY: lint
lint: ## Run linters
	@echo "Running linters..."
	@golangci-lint run --timeout=5m

.PHONY: lint-fix
lint-fix: ## Run linters with auto-fix
	@echo "Running linters with auto-fix..."
	@golangci-lint run --fix --timeout=5m

.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w .
	@goimports -w .

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo "Checking code format..."
	@test -z "$$(gofmt -s -l . | tee /dev/stderr)" || (echo "Please run 'make fmt'" && false)

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

.PHONY: tidy
tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	@go mod tidy

.PHONY: verify
verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

.PHONY: download
download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

.PHONY: examples
examples: ## Run all examples
	@echo "Running examples..."
	@cd examples/basic && go run main.go

.PHONY: security
security: ## Run security scan
	@echo "Running security scan..."
	@gosec -quiet ./...

.PHONY: ci
ci: fmt-check vet lint test coverage-ci security ## Run all CI checks locally
	@echo "✅ All CI checks passed!"

.PHONY: pre-commit
pre-commit: fmt lint test-short ## Run pre-commit checks
	@echo "✅ Pre-commit checks passed!"

.PHONY: build
build: ## Build the project
	@echo "Building..."
	@go build -v ./...

.PHONY: clean
clean: ## Clean build artifacts and cache
	@echo "Cleaning..."
	@go clean -cache -testcache -modcache
	@rm -f coverage.out coverage.html benchmark.txt
	@echo "✅ Cleaned!"

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "✅ Tools installed!"

.PHONY: mod-outdated
mod-outdated: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	@go list -u -m all

.PHONY: update-deps
update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated!"

.PHONY: stats
stats: ## Show code statistics
	@echo "Code statistics:"
	@echo "Total lines:"
	@find . -name '*.go' -not -path "./vendor/*" -not -path "./.git/*" | xargs wc -l | tail -n 1
	@echo "\nFiles by type:"
	@find . -name '*.go' -not -path "./vendor/*" -not -path "./.git/*" | wc -l | awk '{print "Go files: " $$1}'
	@find . -name '*_test.go' -not -path "./vendor/*" | wc -l | awk '{print "Test files: " $$1}'

.PHONY: doc
doc: ## Generate and serve documentation
	@echo "Starting documentation server at http://localhost:6060"
	@godoc -http=:6060

.PHONY: release-check
release-check: ci examples ## Check if ready for release
	@echo "✅ Release checks passed!"
	@echo "Ready to release!"

# Default target
.DEFAULT_GOAL := help