# GitLab CI/CD Pipeline for Go Eloquent ORM

This document describes the GitLab CI/CD pipeline configuration for the Go Eloquent ORM project.

## Pipeline Overview

The pipeline is triggered on:
- **Push to main branch**
- **Merge requests to main branch**

## Pipeline Stages

### 1. Test Stage

#### Primary Testing (SQLite)
- **Job**: `test:sqlite`
- **Purpose**: Core functionality testing using SQLite in-memory database
- **Tests Run**:
  - `./tests/model_test.go` - Model CRUD operations
  - `./tests/integration/integration_test.go` - Integration tests
  - `./querybuilder_test.go` - Query builder functionality
  - `./relationships_test.go` - Model relationships

#### Extended Testing (PostgreSQL)
- **Job**: `test:postgresql`
- **Purpose**: Extended testing with PostgreSQL database
- **Status**: `allow_failure: true` (non-blocking)
- **Database**: PostgreSQL 15
- **Configuration**:
  - Database: `test_db`
  - User: `test_user`
  - Password: `test_password`

#### Extended Testing (MySQL)
- **Job**: `test:mysql`
- **Purpose**: Extended testing with MySQL database
- **Status**: `allow_failure: true` (non-blocking)
- **Database**: MySQL 8.0
- **Configuration**:
  - Database: `test_db`
  - User: `test_user`
  - Password: `test_password`

#### Code Quality
- **Job**: `quality`
- **Purpose**: Code formatting and linting checks
- **Tools**:
  - `go fmt` - Code formatting
  - `go vet` - Static analysis
  - `staticcheck` - Advanced static analysis
  - `goimports` - Import formatting

#### Security Scanning
- **Job**: `security`
- **Purpose**: Security vulnerability scanning
- **Tool**: `gosec` - Go security analyzer
- **Output**: JSON report (`gosec-report.json`)

#### Coverage Report (Main branch only)
- **Job**: `coverage`
- **Purpose**: Generate test coverage reports
- **Minimum Coverage**: 70%
- **Outputs**:
  - `coverage.out` - Coverage profile
  - `coverage.html` - HTML coverage report
  - `coverage.txt` - Text coverage summary

#### Performance Benchmarks (Main branch only)
- **Job**: `benchmark`
- **Purpose**: Run performance benchmarks
- **Output**: `benchmark.txt`
- **Status**: `allow_failure: true`

### 2. Build Stage

#### Build Verification
- **Job**: `build`
- **Purpose**: Verify that the code compiles successfully
- **Dependencies**: Requires `test:sqlite` to pass
- **Command**: `go build -v ./...`

### 3. Deploy Stage

#### Documentation Generation (Main branch only)
- **Job**: `docs`
- **Purpose**: Generate Go documentation
- **Tool**: `godoc`
- **Output**: `docs.txt`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GO_VERSION` | Go version to use | `1.21` |
| `CGO_ENABLED` | Enable CGO (required for SQLite) | `1` |
| `GOCACHE` | Go build cache directory | `.cache/go-build/` |
| `GOMODCACHE` | Go module cache directory | `.cache/go-mod/` |

## Caching

The pipeline uses Go module and build caching to improve performance:
- **Module cache**: `.cache/go-mod/`
- **Build cache**: `.cache/go-build/`

## Artifacts

The pipeline generates several artifacts:

### Test Artifacts
- **Coverage reports** (1 week retention)
- **Security scan reports** (1 week retention)
- **Benchmark results** (1 week retention)

### Documentation Artifacts
- **API documentation** (1 week retention)

## Pipeline Rules

### Triggers
- **Main branch**: All jobs run
- **Merge requests**: All jobs run except coverage and benchmarks
- **Other branches**: Pipeline does not run

### Failure Handling
- **Blocking jobs**: `test:sqlite`, `build`
- **Non-blocking jobs**: `test:postgresql`, `test:mysql`, `quality`, `security`, `benchmark`, `docs`

## Local Development

To run the same tests locally:

```bash
# Run core tests (SQLite)
go test -v ./tests/model_test.go
go test -v ./tests/integration/integration_test.go
go test -v ./querybuilder_test.go
go test -v ./relationships_test.go

# Run code quality checks
go fmt ./...
go vet ./...
staticcheck ./...
goimports -l .

# Generate coverage report
go test -coverprofile=coverage.out ./tests/...
go tool cover -html=coverage.out -o coverage.html

# Run security scan
gosec ./...
```

## Requirements

### System Dependencies
- Go 1.21+
- GCC (for CGO/SQLite)
- Git

### Go Dependencies
- SQLite driver (for primary testing)
- PostgreSQL driver (for extended testing)
- MySQL driver (for extended testing)

## Monitoring

The pipeline provides:
- **Test results**: JUnit XML format
- **Coverage metrics**: Integrated with GitLab coverage display
- **Security reports**: JSON format for security dashboard
- **Performance metrics**: Benchmark results

## Troubleshooting

### Common Issues

1. **CGO compilation errors**: Ensure `CGO_ENABLED=1` and GCC is available
2. **Database connection timeouts**: Database services may need additional startup time
3. **Cache issues**: Clear cache by restarting pipeline with clean cache

### Debug Steps

1. Check job logs for specific error messages
2. Verify Go version and dependencies
3. Ensure database services are properly initialized
4. Check artifact outputs for detailed reports

## Configuration Updates

To modify the pipeline:

1. Edit `.gitlab-ci.yml` in the project root
2. Test changes in a merge request
3. Monitor pipeline execution and adjust as needed

The pipeline is designed to be maintainable and extensible for future testing requirements. 