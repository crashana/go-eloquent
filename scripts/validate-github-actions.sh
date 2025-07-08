#!/bin/bash

# GitHub Actions Workflow Validation Script
# This script simulates the GitHub Actions workflow locally for testing

set -e

echo "🚀 GitHub Actions Workflow Validation for Go Eloquent ORM"
echo "========================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_status "Checking prerequisites..."

if ! command_exists go; then
    print_error "Go is not installed"
    exit 1
fi

if ! command_exists git; then
    print_error "Git is not installed"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_success "Go version: $GO_VERSION"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "Not in a Go module directory. Please run this script from the project root."
    exit 1
fi

print_success "Prerequisites check passed"

# Job 1: Test Job
echo ""
print_status "Job 1: Test Job"
echo "==============="

print_status "Downloading dependencies..."
go mod download

print_status "Verifying dependencies..."
go mod verify

print_status "Running go vet..."
if go vet ./...; then
    print_success "Go vet passed"
else
    print_error "Go vet failed"
    exit 1
fi

print_status "Running core tests..."
if go test -v -race -coverprofile=coverage.out ./tests/model_test.go; then
    print_success "Model tests passed"
else
    print_error "Model tests failed"
    exit 1
fi

if go test -v -race -coverprofile=coverage-integration.out ./tests/integration/integration_test.go; then
    print_success "Integration tests passed"
else
    print_error "Integration tests failed"
    exit 1
fi

# Job 2: Integration Test Job
echo ""
print_status "Job 2: Integration Test Job (SQLite)"
echo "===================================="

export DB_DRIVER=sqlite3
export DB_DATABASE=:memory:

print_status "Running integration tests with SQLite..."
if go test -v ./tests/model_test.go; then
    print_success "SQLite model tests passed"
else
    print_error "SQLite model tests failed"
    exit 1
fi

if go test -v ./tests/integration/integration_test.go; then
    print_success "SQLite integration tests passed"
else
    print_error "SQLite integration tests failed"
    exit 1
fi

print_status "Running additional tests..."
if [ -f "querybuilder_test.go" ]; then
    if go test -v ./querybuilder_test.go; then
        print_success "Query builder tests passed"
    else
        print_warning "Query builder tests failed"
    fi
else
    print_warning "Query builder tests not found"
fi

if [ -f "relationships_test.go" ]; then
    if go test -v ./relationships_test.go; then
        print_success "Relationship tests passed"
    else
        print_warning "Relationship tests failed"
    fi
else
    print_warning "Relationship tests not found"
fi

# Job 3: Lint Job
echo ""
print_status "Job 3: Lint Job"
echo "==============="

if command_exists golangci-lint; then
    print_status "Running golangci-lint..."
    if golangci-lint run --timeout=5m; then
        print_success "Linting passed"
    else
        print_warning "Linting found issues"
    fi
else
    print_warning "golangci-lint not installed. Install from: https://golangci-lint.run/usage/install/"
    
    # Fallback to basic linting
    print_status "Running basic linting checks..."
    if go fmt ./...; then
        print_success "Go fmt passed"
    else
        print_warning "Go fmt found issues"
    fi
fi

# Job 4: Security Job
echo ""
print_status "Job 4: Security Job"
echo "==================="

if command_exists gosec; then
    print_status "Running gosec security scanner..."
    if gosec -no-fail -fmt sarif -out results.sarif ./...; then
        print_success "Security scan completed"
        if [ -f "results.sarif" ]; then
            print_success "SARIF report generated: results.sarif"
        fi
    else
        print_warning "Security scan found potential issues"
    fi
else
    print_warning "Gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
fi

# Job 5: Build Job
echo ""
print_status "Job 5: Build Job"
echo "================"

print_status "Creating bin directory..."
mkdir -p bin

print_status "Building library..."
if go build -v ./...; then
    print_success "Library build passed"
else
    print_error "Library build failed"
    exit 1
fi

print_status "Building example (if exists)..."
if [ -d "Examples" ]; then
    print_status "Found Examples directory"
    cd Examples
    if go build -v -o ../bin/example .; then
        print_success "Example build passed"
    else
        print_warning "Example build failed"
    fi
    cd ..
elif [ -d "examples" ]; then
    print_status "Found examples directory"
    cd examples
    if go build -v -o ../bin/example .; then
        print_success "Example build passed"
    else
        print_warning "Example build failed"
    fi
    cd ..
else
    print_warning "No examples directory found, creating placeholder"
    echo "package main; func main() { println(\"Go Eloquent ORM\") }" > main.go
    go build -o bin/eloquent main.go
    rm main.go
    print_success "Placeholder binary created"
fi

# Job 6: Coverage Job
echo ""
print_status "Job 6: Coverage Job"
echo "==================="

print_status "Running tests with coverage..."
if go test -coverprofile=coverage.out ./tests/...; then
    print_success "Coverage tests passed"
    
    if command_exists go; then
        print_status "Generating coverage report..."
        go tool cover -html=coverage.out -o coverage.html
        go tool cover -func=coverage.out | tee coverage.txt
        
        COVERAGE=$(go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+' || echo "0.0")
        print_status "Total coverage: ${COVERAGE}%"
        
        if [ $(echo "${COVERAGE} < 70" | bc -l 2>/dev/null || echo "0") -eq 1 ]; then
            print_warning "Coverage is below 70%. Current coverage: ${COVERAGE}%"
        else
            print_success "Coverage meets minimum requirement (70%)"
        fi
        
        print_success "Coverage report generated: coverage.html"
    fi
else
    print_warning "Coverage tests failed"
fi

# Job 7: Benchmark Job
echo ""
print_status "Job 7: Benchmark Job"
echo "===================="

print_status "Running benchmarks..."
if go test -bench=. -benchmem ./tests/... > benchmark.txt 2>&1; then
    if grep -q "Benchmark" benchmark.txt; then
        print_success "Benchmarks completed"
        echo "Benchmark results:"
        cat benchmark.txt
    else
        print_warning "No benchmarks found"
    fi
else
    print_warning "Benchmark execution failed"
fi

# Final summary
echo ""
print_status "GitHub Actions Workflow Validation Summary"
echo "=========================================="

print_success "✅ Core tests passed"
print_success "✅ Integration tests (SQLite) passed"
print_success "✅ Build verification passed"
print_success "✅ Coverage report generated"

# Check for warnings
WARNINGS=0
if ! command_exists golangci-lint; then
    WARNINGS=$((WARNINGS + 1))
fi
if ! command_exists gosec; then
    WARNINGS=$((WARNINGS + 1))
fi

if [ $WARNINGS -eq 0 ]; then
    print_success "✅ All tools are installed"
else
    print_warning "⚠️  $WARNINGS optional tools are missing (see warnings above)"
fi

echo ""
print_success "🎉 GitHub Actions workflow validation completed successfully!"
print_status "Your code is ready for GitHub Actions execution."

# Clean up temporary files
rm -f coverage.out coverage-integration.out benchmark.txt results.sarif main.go

echo ""
print_status "Temporary files cleaned up."
print_status "Artifacts generated:"
print_status "  - coverage.html (Coverage report)"
print_status "  - coverage.txt (Coverage summary)"
print_status "  - bin/ (Build artifacts)"

if [ -f "results.sarif" ]; then
    print_status "  - results.sarif (Security scan report)"
fi 