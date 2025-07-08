# GitHub Actions CI/CD Pipeline for Go Eloquent ORM

This document provides comprehensive information about the GitHub Actions CI/CD pipeline configuration for the Go Eloquent ORM project.

## Overview

The GitHub Actions workflow provides automated testing, building, and quality assurance for the Go Eloquent ORM package. It runs on every push to `main` and `develop` branches, as well as on pull requests targeting these branches.

## Workflow Structure

The pipeline consists of 7 main jobs that run in parallel where possible:

### 1. Test Job
- **Purpose**: Core functionality testing
- **Matrix**: Go versions 1.19, 1.20, 1.21
- **Steps**:
  - Checkout code
  - Set up Go environment
  - Cache Go modules
  - Download and verify dependencies
  - Run `go vet` for static analysis
  - Execute core tests (model and integration tests)
  - Upload coverage to Codecov

### 2. Integration Test Job
- **Purpose**: Database integration testing
- **Matrix**: 
  - Go versions: 1.20, 1.21
  - Databases: SQLite, PostgreSQL, MySQL
- **Services**: PostgreSQL 15, MySQL 8.0
- **Steps**:
  - Set up test environment with database configurations
  - Wait for database services to be ready
  - Run integration tests with different database backends
  - Run additional tests (query builder, relationships)

### 3. Lint Job
- **Purpose**: Code quality and style checking
- **Tools**: golangci-lint
- **Steps**:
  - Run comprehensive linting with 5-minute timeout
  - Check code formatting, style, and potential issues

### 4. Security Job
- **Purpose**: Security vulnerability scanning
- **Tools**: Gosec
- **Steps**:
  - Run security scanner
  - Generate SARIF report
  - Upload results to GitHub Security tab

### 5. Build Job
- **Purpose**: Cross-platform build verification
- **Matrix**: 
  - OS: Linux, Windows, macOS
  - Architecture: amd64, arm64
  - Excludes: Windows/arm64
- **Steps**:
  - Build library for all platforms
  - Build example application (if exists)
  - Upload build artifacts

### 6. Coverage Job
- **Purpose**: Code coverage analysis
- **Dependencies**: Requires test job to complete
- **Steps**:
  - Run tests with coverage profiling
  - Generate HTML coverage report
  - Upload coverage artifacts

### 7. Benchmark Job
- **Purpose**: Performance benchmarking
- **Trigger**: Only on main branch
- **Steps**:
  - Run benchmark tests
  - Upload benchmark results

## Key Features

### ✅ Fixed Issues from Previous Configuration

1. **Deprecated Actions**: Updated to latest versions
   - `actions/cache@v4`
   - `actions/upload-artifact@v4`
   - `codecov/codecov-action@v4`
   - `golangci/golangci-lint-action@v4`

2. **Go Version Matrix**: Removed invalid Go 1.2, using 1.19, 1.20, 1.21

3. **Test Execution**: Fixed test commands to use specific test files

4. **Database Configuration**: Fixed environment variable names (`DB_DRIVER` instead of `DB_CONNECTION`)

5. **Service Readiness**: Added proper service waiting logic

6. **Build Robustness**: Added error handling for missing Examples directory

7. **Artifact Management**: Improved artifact handling with proper naming and conditional uploads

### 🚀 Enhanced Features

- **Parallel Execution**: Jobs run in parallel for faster CI/CD
- **Smart Caching**: Go modules and build cache for faster builds
- **Cross-Platform Support**: Build verification for multiple OS/architecture combinations
- **Comprehensive Testing**: Unit, integration, and database-specific tests
- **Quality Gates**: Linting, security scanning, and coverage reporting
- **Artifact Management**: Build artifacts, coverage reports, and benchmark results
- **Error Resilience**: Graceful handling of optional components

## Environment Variables

### Database Configuration

The pipeline automatically configures these environment variables based on the database matrix:

#### SQLite
```bash
DB_DRIVER=sqlite3
DB_DATABASE=:memory:
```

#### PostgreSQL
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=eloquent_test
DB_USERNAME=postgres
DB_PASSWORD=postgres
```

#### MySQL
```bash
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=eloquent_test
DB_USERNAME=root
DB_PASSWORD=root
```

### Required Secrets

Add these secrets to your GitHub repository:

1. **CODECOV_TOKEN** (optional): For coverage reporting
   - Go to [Codecov](https://codecov.io/)
   - Add your repository
   - Copy the token to GitHub Secrets

## Local Validation

Use the provided validation script to test the workflow locally:

```bash
# Make script executable
chmod +x scripts/validate-github-actions.sh

# Run validation
./scripts/validate-github-actions.sh
```

The script will:
- Check prerequisites
- Run all test jobs locally
- Generate coverage reports
- Verify build process
- Provide colored output with status information

## Artifacts Generated

The workflow generates several artifacts:

1. **Build Artifacts**: Cross-platform binaries
2. **Coverage Reports**: HTML and text coverage reports
3. **Benchmark Results**: Performance benchmark data
4. **Security Reports**: SARIF format security scan results

## Monitoring and Troubleshooting

### Common Issues

1. **Test Failures**: Check test logs for specific error messages
2. **Build Failures**: Verify Go module dependencies
3. **Coverage Issues**: Ensure tests are properly structured
4. **Lint Failures**: Run `golangci-lint` locally to fix issues

### Debugging Tips

1. **Local Testing**: Always run `./scripts/validate-github-actions.sh` before pushing
2. **Dependency Issues**: Check `go.mod` and `go.sum` files
3. **Database Tests**: Verify database connection strings and schemas
4. **Cross-Platform**: Test builds on different platforms if possible

### Performance Optimization

The workflow includes several optimizations:

- **Caching**: Go modules and build cache
- **Parallel Jobs**: Independent jobs run simultaneously
- **Conditional Execution**: Some jobs only run on specific branches
- **Smart Artifacts**: Only upload when files exist

## Workflow Triggers

The workflow runs on:

- **Push** to `main` or `develop` branches
- **Pull Request** targeting `main` or `develop` branches

## Required Dependencies

### Go Modules
The project requires these Go modules (automatically handled):
- Database drivers (sqlite3, postgres, mysql)
- Testing frameworks
- Any project-specific dependencies

### External Tools (Optional)
- **golangci-lint**: For advanced linting
- **gosec**: For security scanning
- **codecov**: For coverage reporting

## Best Practices

1. **Test Locally**: Always run validation script before pushing
2. **Keep Dependencies Updated**: Regularly update Go modules
3. **Monitor Coverage**: Maintain good test coverage
4. **Security Scanning**: Address security issues promptly
5. **Performance Monitoring**: Review benchmark results regularly

## Customization

### Adding New Go Versions
Update the matrix in `.github/workflows/test.yml`:

```yaml
strategy:
  matrix:
    go-version: [1.19, 1.20, 1.21, 1.22]  # Add new versions
```

### Adding New Database Support
1. Add database to the matrix
2. Add service configuration
3. Update environment variable setup
4. Add database-specific test configuration

### Modifying Test Commands
Update the test execution steps to match your project structure:

```yaml
- name: Run tests
  run: |
    go test -v ./your/test/path
    go test -v ./another/test/path
```

## Support

For issues related to:
- **GitHub Actions**: Check GitHub Actions documentation
- **Go Testing**: Refer to Go testing documentation
- **Database Integration**: Check database-specific documentation
- **Project-Specific Issues**: Create an issue in the project repository

## Changelog

### v1.0.0 (Current)
- Initial GitHub Actions workflow
- Multi-database testing support
- Cross-platform build verification
- Security scanning integration
- Coverage reporting
- Benchmark testing
- Local validation script

---

This documentation is maintained alongside the GitHub Actions workflow. Please update it when making changes to the CI/CD configuration. 