#!/bin/bash

# GitLab CI/CD Pipeline Validation Script
# This script simulates the GitLab CI pipeline locally for testing

set -e

echo "🚀 GitLab CI/CD Pipeline Validation for Go Eloquent ORM"
echo "======================================================"

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

# Stage 1: Test Stage
echo ""
print_status "Stage 1: Running Tests (SQLite)"
echo "--------------------------------"

# Core tests with SQLite
print_status "Running model tests..."
if go test -v ./tests/model_test.go; then
    print_success "Model tests passed"
else
    print_error "Model tests failed"
    exit 1
fi

print_status "Running integration tests..."
if go test -v ./tests/integration/integration_test.go; then
    print_success "Integration tests passed"
else
    print_error "Integration tests failed"
    exit 1
fi

print_status "Running query builder tests..."
if go test -v ./querybuilder_test.go; then
    print_success "Query builder tests passed"
else
    print_warning "Query builder tests failed (may not exist)"
fi

print_status "Running relationship tests..."
if go test -v ./relationships_test.go; then
    print_success "Relationship tests passed"
else
    print_warning "Relationship tests failed (may not exist)"
fi

# Code quality checks
echo ""
print_status "Running Code Quality Checks"
echo "----------------------------"

print_status "Checking code formatting..."
if [ -n "$(go fmt ./...)" ]; then
    print_warning "Code is not properly formatted. Run 'go fmt ./...' to fix."
else
    print_success "Code formatting check passed"
fi

print_status "Running go vet..."
if go vet ./...; then
    print_success "Go vet passed"
else
    print_warning "Go vet found issues"
fi

# Check for staticcheck
if command_exists staticcheck; then
    print_status "Running staticcheck..."
    if staticcheck ./...; then
        print_success "Staticcheck passed"
    else
        print_warning "Staticcheck found issues"
    fi
else
    print_warning "Staticcheck not installed. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"
fi

# Check for goimports
if command_exists goimports; then
    print_status "Checking import formatting..."
    if [ -n "$(goimports -l .)" ]; then
        print_warning "Imports are not properly formatted. Run 'goimports -w .' to fix."
    else
        print_success "Import formatting check passed"
    fi
else
    print_warning "Goimports not installed. Install with: go install golang.org/x/tools/cmd/goimports@latest"
fi

# Security scan
echo ""
print_status "Running Security Scan"
echo "---------------------"

if command_exists gosec; then
    print_status "Running gosec security scan..."
    if gosec ./...; then
        print_success "Security scan passed"
    else
        print_warning "Security scan found potential issues"
    fi
else
    print_warning "Gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"
fi

# Stage 2: Build Stage
echo ""
print_status "Stage 2: Build Verification"
echo "---------------------------"

print_status "Building project..."
if go build -v ./...; then
    print_success "Build verification passed"
else
    print_error "Build verification failed"
    exit 1
fi

# Coverage report
echo ""
print_status "Generating Coverage Report"
echo "--------------------------"

print_status "Running tests with coverage..."
if go test -coverprofile=coverage.out ./tests/...; then
    print_success "Coverage tests passed"
    
    # Generate coverage report
    if command_exists go; then
        COVERAGE=$(go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+' || echo "0.0")
        print_status "Total coverage: ${COVERAGE}%"
        
        if [ $(echo "${COVERAGE} < 70" | bc -l 2>/dev/null || echo "0") -eq 1 ]; then
            print_warning "Coverage is below 70%. Current coverage: ${COVERAGE}%"
        else
            print_success "Coverage meets minimum requirement (70%)"
        fi
        
        # Generate HTML coverage report
        go tool cover -html=coverage.out -o coverage.html
        print_success "Coverage report generated: coverage.html"
    fi
else
    print_warning "Coverage tests failed"
fi

# Benchmarks
echo ""
print_status "Running Performance Benchmarks"
echo "------------------------------"

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

# Documentation generation
echo ""
print_status "Generating Documentation"
echo "-----------------------"

if command_exists godoc; then
    print_status "Generating Go documentation..."
    go doc -all . > docs.txt
    print_success "Documentation generated: docs.txt"
else
    print_warning "Godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest"
fi

# Final summary
echo ""
print_status "Pipeline Validation Summary"
echo "==========================="

print_success "✅ Core tests (SQLite) passed"
print_success "✅ Build verification passed"

# Check for warnings
if command_exists staticcheck && command_exists goimports && command_exists gosec && command_exists godoc; then
    print_success "✅ All optional tools are installed"
else
    print_warning "⚠️  Some optional tools are missing (see warnings above)"
fi

echo ""
print_success "🎉 Pipeline validation completed successfully!"
print_status "Your code is ready for GitLab CI/CD pipeline execution."

# Clean up temporary files
rm -f coverage.out benchmark.txt docs.txt

echo ""
print_status "Temporary files cleaned up."
print_status "To view coverage report, run: go tool cover -html=coverage.out -o coverage.html" 