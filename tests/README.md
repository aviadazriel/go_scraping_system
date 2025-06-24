# Testing Guide

This document provides comprehensive information about the testing strategy, test organization, and how to run tests for the Go Scraping Project.

## Test Organization

The project follows a structured testing approach with the following organization:

```
tests/
├── unit/                    # Unit tests
│   ├── url_repository_test.go
│   ├── url_scheduler_test.go
│   └── ...
├── integration/             # Integration tests
│   ├── api_gateway_test.go
│   ├── database_test.go
│   └── ...
├── utils/                   # Test utilities and helpers
│   └── test_helpers.go
├── fixtures/                # Test data and fixtures
│   ├── sample_urls.json
│   └── ...
└── README.md               # This file
```

## Test Types

### 1. Unit Tests (`tests/unit/`)

Unit tests focus on testing individual functions and methods in isolation. They use mocks to isolate the code under test from external dependencies.

**Characteristics:**
- Fast execution (< 1 second per test)
- No external dependencies (database, network, etc.)
- Use mocks for external interfaces
- Test business logic in isolation

**Example:**
```go
func TestURLRepository_CreateURL(t *testing.T) {
    // Test URL repository with mocked database
    mockQuerier := &MockQuerier{}
    repo := repositories.NewURLRepository(mockQuerier, logrus.New())
    
    // Test business logic without real database
    result, err := repo.CreateURL(context.Background(), request)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

### 2. Integration Tests (`tests/integration/`)

Integration tests verify that different components work together correctly. They may use real databases, HTTP servers, or other external services.

**Characteristics:**
- Slower execution (1-10 seconds per test)
- May use real external dependencies
- Test component interactions
- Verify end-to-end functionality

**Example:**
```go
func TestAPIGateway_CreateURL_Integration(t *testing.T) {
    // Setup real HTTP server with mocked dependencies
    server, mockRepo := setupTestServer()
    defer server.Close()
    
    // Make real HTTP request
    resp, err := http.Post(server.URL+"/api/v1/urls", "application/json", body)
    
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
```

### 3. Test Utilities (`tests/utils/`)

Test utilities provide common functionality used across different test types:

- Database setup and cleanup
- Mock creation helpers
- Test data generators
- Assertion helpers
- Environment configuration

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make test-benchmark
```

### Advanced Test Commands

```bash
# Run tests in parallel
make test-parallel

# Run tests with verbose output
make test-verbose

# Run tests for specific services
make test-url-manager
make test-api-gateway
make test-database

# Run performance tests
make test-performance

# Run security tests
make test-security

# Run full test suite (setup, tests, cleanup)
make test-full
```

### Test Environment Setup

```bash
# Setup test environment
make test-setup

# Setup test database
make test-db-setup

# Cleanup test artifacts
make test-clean

# Cleanup test database
make test-db-cleanup
```

### CI/CD Test Commands

```bash
# Run CI/CD pipeline simulation
make test-ci

# Test with different Go versions
make test-go-versions

# Generate coverage reports
make test-coverage-badge
```

## Test Configuration

### Environment Variables

Tests can be configured using environment variables:

```bash
# Database configuration
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=scraper
export TEST_DB_PASSWORD=scraper_password
export TEST_DB_NAME=scraper_test
export TEST_DB_SSLMODE=disable

# Kafka configuration
export TEST_KAFKA_BROKERS=localhost:9092
export TEST_KAFKA_GROUP_ID=test-group

# Test configuration
export TEST_TIMEOUT=30s
export TEST_PARALLEL=4
```

### Test Tags

Tests can be tagged for selective execution:

```bash
# Run short tests only
go test -v -short ./tests/...

# Run long tests only
go test -v -run "Test.*Long" ./tests/...

# Run tests with specific tags
go test -v -tags=race ./tests/...
go test -v -tags=debug ./tests/...
```

## Writing Tests

### Unit Test Guidelines

1. **Use table-driven tests** for multiple test cases:
```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        expectError bool
    }{
        {"valid input", "test", false},
        {"invalid input", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Function(tt.input)
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, result)
            }
        })
    }
}
```

2. **Use descriptive test names** that explain the scenario:
```go
func TestURLRepository_CreateURL_WithValidData_ReturnsSuccess(t *testing.T)
func TestURLRepository_CreateURL_WithInvalidURL_ReturnsError(t *testing.T)
```

3. **Mock external dependencies**:
```go
mockRepo := &MockURLRepository{}
mockRepo.On("CreateURL", mock.Anything, expectedRequest).Return(expectedResponse, nil)
```

### Integration Test Guidelines

1. **Setup and teardown** properly:
```go
func TestIntegration(t *testing.T) {
    // Setup
    server, cleanup := setupTestServer(t)
    defer cleanup()
    
    // Test logic
    // ...
}
```

2. **Use real HTTP clients** for API testing:
```go
resp, err := http.Post(server.URL+"/api/v1/urls", "application/json", body)
assert.NoError(t, err)
assert.Equal(t, http.StatusCreated, resp.StatusCode)
```

3. **Test error scenarios**:
```go
// Test invalid request
resp, err := http.Post(server.URL+"/api/v1/urls", "application/json", invalidBody)
assert.NoError(t, err)
assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
```

### Test Utilities Usage

Use the provided test utilities for common operations:

```go
import "go_scraping_project/tests/utils"

func TestExample(t *testing.T) {
    // Setup test database
    db, cleanup := utils.SetupTestDatabase(t)
    defer cleanup()
    
    // Create test data
    url := utils.CreateTestURL(t, db, "https://example.com", "hourly", "Test URL")
    
    // Cleanup test data
    utils.CleanupTestData(t, db)
    
    // Use test helpers
    utils.ExpectNoError(t, err)
    utils.AssertHTTPStatus(t, resp.StatusCode, http.StatusOK)
}
```

## Test Coverage

### Coverage Targets

- **Unit tests**: Aim for >90% coverage
- **Integration tests**: Aim for >80% coverage
- **Overall coverage**: Aim for >85% coverage

### Coverage Reports

Generate coverage reports:

```bash
# Generate HTML coverage report
make test-coverage

# Generate XML coverage report for CI
make test-coverage-badge
```

View coverage reports:
- HTML: Open `coverage.html` in a browser
- XML: Use with CI/CD tools like Codecov

## Best Practices

### 1. Test Organization

- Keep tests close to the code they test
- Use descriptive test names
- Group related tests together
- Use subtests for multiple scenarios

### 2. Test Data

- Use fixtures for complex test data
- Generate unique test data to avoid conflicts
- Clean up test data after tests
- Use constants for common test values

### 3. Mocking

- Mock external dependencies
- Use interfaces for testability
- Verify mock expectations
- Keep mocks simple and focused

### 4. Assertions

- Use descriptive assertion messages
- Test both positive and negative cases
- Verify error conditions
- Check edge cases

### 5. Performance

- Keep tests fast
- Use parallel execution where possible
- Avoid unnecessary setup/teardown
- Use test timeouts for long-running tests

## Troubleshooting

### Common Issues

1. **Test database connection failed**:
   - Check if PostgreSQL is running
   - Verify database credentials
   - Ensure test database exists

2. **Tests failing intermittently**:
   - Check for race conditions
   - Use `-race` flag to detect races
   - Ensure proper cleanup

3. **Slow test execution**:
   - Use parallel execution
   - Optimize database queries
   - Reduce unnecessary setup

4. **Mock expectations not met**:
   - Check mock setup
   - Verify method calls
   - Use `mock.Anything` for flexible matching

### Debugging Tests

```bash
# Run tests with verbose output
go test -v ./tests/...

# Run specific test with debug output
go test -v -run TestSpecificFunction ./tests/...

# Run tests with race detection
go test -race ./tests/...

# Run tests with coverage and show uncovered lines
go test -coverprofile=coverage.out ./tests/...
go tool cover -func=coverage.out
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make test-ci
      - run: make test-coverage-badge
      - uses: codecov/codecov-action@v3
        with:
          file: ./coverage.xml
```

### Local CI Simulation

```bash
# Run full CI pipeline locally
make test-ci
```

## Contributing

When adding new tests:

1. Follow the existing test patterns
2. Add tests for new functionality
3. Update this documentation if needed
4. Ensure tests pass in CI
5. Maintain good test coverage

## Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Testify Assertion Library](https://github.com/stretchr/testify)
- [Go Test Examples](https://golang.org/doc/code.html#Testing)
- [Testing Best Practices](https://github.com/golang/go/wiki/CodeReviewComments#tests) 