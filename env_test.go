package eloquent

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	// Create a temporary .env file
	envContent := `# Test environment file
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=test_db
DB_USERNAME=test_user
DB_PASSWORD=test_password
DB_CONNECTION=pgsql
DEBUG=true
EMPTY_VALUE=
QUOTED_VALUE="quoted string"
SINGLE_QUOTED='single quoted'
`

	// Write to temporary file
	tmpFile := "test.env"
	err := os.WriteFile(tmpFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test env file: %v", err)
	}
	defer os.Remove(tmpFile)

	// Test loading the env file
	err = LoadEnv(tmpFile)
	if err != nil {
		t.Fatalf("LoadEnv failed: %v", err)
	}

	// Test that values were loaded correctly using Env function
	tests := []struct {
		key      string
		expected string
	}{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"DB_DATABASE", "test_db"},
		{"DB_USERNAME", "test_user"},
		{"DB_PASSWORD", "test_password"},
		{"DB_CONNECTION", "pgsql"},
		{"DEBUG", "true"},
		{"EMPTY_VALUE", ""},
		{"QUOTED_VALUE", "quoted string"},
		{"SINGLE_QUOTED", "single quoted"},
	}

	for _, test := range tests {
		actual := Env(test.key)
		if actual != test.expected {
			t.Errorf("Expected %s=%s, got %s", test.key, test.expected, actual)
		}
	}
}

func TestLoadEnvNonExistentFile(t *testing.T) {
	err := LoadEnv("nonexistent.env")
	// LoadEnv should not return an error for non-existent files
	// It should just use system environment variables
	if err != nil {
		t.Errorf("Expected no error for non-existent file, got %v", err)
	}
}

func TestEnv(t *testing.T) {
	// Set up test environment variables
	os.Setenv("TEST_VAR", "test_value")
	os.Setenv("EMPTY_VAR", "")
	defer func() {
		os.Unsetenv("TEST_VAR")
		os.Unsetenv("EMPTY_VAR")
	}()

	tests := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{"TEST_VAR", "default", "test_value"},
		{"NONEXISTENT_VAR", "default", "default"},
		{"EMPTY_VAR", "default", "default"}, // Empty env var should return default
	}

	for _, test := range tests {
		actual := Env(test.key, test.defaultValue)
		if actual != test.expected {
			t.Errorf("Env(%s, %s) = %s, expected %s", test.key, test.defaultValue, actual, test.expected)
		}
	}
}

func TestEnvInt(t *testing.T) {
	// Set up test environment variables
	os.Setenv("INT_VAR", "42")
	os.Setenv("INVALID_INT", "not_a_number")
	os.Setenv("EMPTY_INT", "")
	defer func() {
		os.Unsetenv("INT_VAR")
		os.Unsetenv("INVALID_INT")
		os.Unsetenv("EMPTY_INT")
	}()

	tests := []struct {
		key          string
		defaultValue int
		expected     int
	}{
		{"INT_VAR", 0, 42},
		{"NONEXISTENT_INT", 100, 100},
		{"INVALID_INT", 50, 50},
		{"EMPTY_INT", 25, 25},
	}

	for _, test := range tests {
		actual := EnvInt(test.key, test.defaultValue)
		if actual != test.expected {
			t.Errorf("EnvInt(%s, %d) = %d, expected %d", test.key, test.defaultValue, actual, test.expected)
		}
	}
}

func TestEnvBool(t *testing.T) {
	// Set up test environment variables
	os.Setenv("BOOL_TRUE", "true")
	os.Setenv("BOOL_FALSE", "false")
	os.Setenv("BOOL_1", "1")
	os.Setenv("BOOL_0", "0")
	os.Setenv("BOOL_YES", "yes")
	os.Setenv("BOOL_NO", "no")
	os.Setenv("BOOL_INVALID", "invalid")
	os.Setenv("BOOL_EMPTY", "")
	defer func() {
		os.Unsetenv("BOOL_TRUE")
		os.Unsetenv("BOOL_FALSE")
		os.Unsetenv("BOOL_1")
		os.Unsetenv("BOOL_0")
		os.Unsetenv("BOOL_YES")
		os.Unsetenv("BOOL_NO")
		os.Unsetenv("BOOL_INVALID")
		os.Unsetenv("BOOL_EMPTY")
	}()

	tests := []struct {
		key          string
		defaultValue bool
		expected     bool
	}{
		{"BOOL_TRUE", false, true},
		{"BOOL_FALSE", true, false},
		{"BOOL_1", false, true},
		{"BOOL_0", true, false},
		{"BOOL_YES", false, true},
		{"BOOL_NO", true, false},
		{"BOOL_INVALID", true, false}, // Invalid bool should return false, not default
		{"BOOL_EMPTY", true, true},
		{"NONEXISTENT_BOOL", false, false},
	}

	for _, test := range tests {
		actual := EnvBool(test.key, test.defaultValue)
		if actual != test.expected {
			t.Errorf("EnvBool(%s, %t) = %t, expected %t", test.key, test.defaultValue, actual, test.expected)
		}
	}
}

func TestAutoConnect(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{"DB_CONNECTION", "DB_HOST", "DB_PORT", "DB_DATABASE", "DB_USERNAME", "DB_PASSWORD"}
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Clean up after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Test case 1: Missing required fields
	for _, key := range envVars {
		os.Unsetenv(key)
	}

	err := AutoConnect()
	if err == nil {
		t.Error("Expected error for missing required fields, got nil")
	}

	// Test case 2: Valid PostgreSQL configuration
	os.Setenv("DB_CONNECTION", "pgsql")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_DATABASE", "test_db")
	os.Setenv("DB_USERNAME", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")

	// This will fail because we don't have a real database, but it should get past validation
	err = AutoConnect()
	// We expect this to fail with a connection error, not a validation error
	if err != nil && err.Error() == "DB_DATABASE is required" {
		t.Error("AutoConnect should pass validation with valid config")
	}

	// Test case 3: Valid MySQL configuration
	os.Setenv("DB_CONNECTION", "mysql")
	os.Setenv("DB_PORT", "3306")

	err = AutoConnect()
	// We expect this to fail with a connection error, not a validation error
	if err != nil && err.Error() == "DB_DATABASE is required" {
		t.Error("AutoConnect should pass validation with valid MySQL config")
	}
}

func TestParseEnvLine(t *testing.T) {
	tests := []struct {
		line     string
		expected map[string]string
	}{
		{"KEY=value", map[string]string{"KEY": "value"}},
		{"KEY=", map[string]string{"KEY": ""}},
		{"KEY=\"quoted value\"", map[string]string{"KEY": "quoted value"}},
		{"KEY='single quoted'", map[string]string{"KEY": "single quoted"}},
		{"# This is a comment", map[string]string{}},
		{"", map[string]string{}},
		{"   ", map[string]string{}},
		{"INVALID_LINE", map[string]string{}},
		{"KEY=value # with comment", map[string]string{"KEY": "value"}},
		{"KEY=value with spaces", map[string]string{"KEY": "value with spaces"}},
	}

	for _, test := range tests {
		// Create temporary file with single line
		tmpFile := "single_line.env"
		err := os.WriteFile(tmpFile, []byte(test.line), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer os.Remove(tmpFile)

		// Load the env file
		err = LoadEnv(tmpFile)
		if err != nil {
			t.Fatalf("Failed to load env file: %v", err)
		}

		// Check results using Env function
		for key, expectedValue := range test.expected {
			actualValue := Env(key)
			if actualValue != expectedValue {
				t.Errorf("Line '%s': expected %s=%s, got %s", test.line, key, expectedValue, actualValue)
			}
		}

		// Reset envConfig for next test
		envConfig = nil
	}
}
