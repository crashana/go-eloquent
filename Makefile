.PHONY: test test-unit test-integration test-coverage clean build lint fmt vet deps help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Binary names
BINARY_NAME=eloquent-example
BINARY_UNIX=$(BINARY_NAME)_unix

# Test parameters
TEST_PACKAGES=./...
TEST_TIMEOUT=30s
COVERAGE_FILE=coverage.out

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests
	$(GOTEST) -v -race -timeout $(TEST_TIMEOUT) ./env_test.go ./connection_test.go ./querybuilder_test.go ./relationships_test.go

test-integration: ## Run integration tests
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./tests/...

test-models: ## Run model tests
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./tests/model_test.go

test-coverage: ## Run tests with coverage
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(TEST_PACKAGES)
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html

test-coverage-ci: ## Run tests with coverage for CI
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic $(TEST_PACKAGES)

lint: ## Run linter
	golangci-lint run --timeout=5m

fmt: ## Format code
	$(GOFMT) -s -w .

vet: ## Run go vet
	$(GOVET) $(TEST_PACKAGES)

build: ## Build the example application
	cd Examples && $(GOBUILD) -o ../bin/$(BINARY_NAME) -v .

build-linux: ## Build for Linux
	cd Examples && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ../bin/$(BINARY_UNIX) -v .

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)
	rm -f $(COVERAGE_FILE)
	rm -f coverage.html

run: ## Run the example application
	cd Examples && $(GOCMD) run main.go

setup-db: ## Set up test database (SQLite)
	@echo "Setting up SQLite test database..."
	@echo "DB_CONNECTION=sqlite3" > .env.test
	@echo "DB_DATABASE=:memory:" >> .env.test
	@echo "Test environment configured for SQLite"

setup-postgres: ## Set up PostgreSQL test environment
	@echo "Setting up PostgreSQL test environment..."
	@echo "DB_CONNECTION=postgres" > .env.test
	@echo "DB_HOST=localhost" >> .env.test
	@echo "DB_PORT=5432" >> .env.test
	@echo "DB_DATABASE=eloquent_test" >> .env.test
	@echo "DB_USERNAME=postgres" >> .env.test
	@echo "DB_PASSWORD=postgres" >> .env.test
	@echo "Test environment configured for PostgreSQL"

setup-mysql: ## Set up MySQL test environment
	@echo "Setting up MySQL test environment..."
	@echo "DB_CONNECTION=mysql" > .env.test
	@echo "DB_HOST=localhost" >> .env.test
	@echo "DB_PORT=3306" >> .env.test
	@echo "DB_DATABASE=eloquent_test" >> .env.test
	@echo "DB_USERNAME=root" >> .env.test
	@echo "DB_PASSWORD=root" >> .env.test
	@echo "Test environment configured for MySQL"

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

benchmark: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem $(TEST_PACKAGES)

check: deps vet lint test ## Run all checks (deps, vet, lint, test)

ci: deps vet test-coverage-ci ## Run CI pipeline

validate-github: ## Validate GitHub Actions workflow locally
	@echo "Validating GitHub Actions workflow..."
	@chmod +x scripts/validate-github-actions.sh
	@./scripts/validate-github-actions.sh

validate-gitlab: ## Validate GitLab CI pipeline locally
	@echo "Validating GitLab CI pipeline..."
	@chmod +x scripts/validate-ci.sh
	@./scripts/validate-ci.sh

docker-test: ## Run tests in Docker
	docker run --rm -v $(PWD):/app -w /app golang:1.21 make test

docker-build: ## Build in Docker
	docker run --rm -v $(PWD):/app -w /app golang:1.21 make build

# Development helpers
watch-test: ## Watch for file changes and run tests
	@echo "Watching for changes... (requires 'entr' command)"
	find . -name "*.go" | entr -c make test-unit

example: ## Run the example with sample data
	cd Examples && $(GOCMD) run main.go

docs: ## Generate documentation
	$(GOCMD) doc -all . > docs/api.md

# Git hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@echo '#!/bin/sh\nmake fmt vet lint' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git hooks installed"

# Release helpers
tag: ## Create a new git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make tag VERSION=v1.0.0"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

release: clean check build ## Prepare for release 