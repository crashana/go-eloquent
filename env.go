package eloquent

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvConfig holds environment configuration
type EnvConfig struct {
	values map[string]string
}

var envConfig *EnvConfig

// LoadEnv loads environment variables from .env file
func LoadEnv(filepath ...string) error {
	envFile := ".env"
	if len(filepath) > 0 {
		envFile = filepath[0]
	}

	config := &EnvConfig{
		values: make(map[string]string),
	}

	// Check if .env file exists
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		// .env file doesn't exist, that's okay - we'll use system environment variables
		envConfig = config
		return nil
	}

	file, err := os.Open(envFile)
	if err != nil {
		return fmt.Errorf("failed to open .env file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Handle comments in values - remove everything after # that's not in quotes
		if commentIndex := strings.Index(value, "#"); commentIndex != -1 {
			// Check if the # is inside quotes
			if !(strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) &&
				!(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
				value = strings.TrimSpace(value[:commentIndex])
			}
		}

		// Remove quotes if present
		if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) ||
			(strings.HasPrefix(value, `'`) && strings.HasSuffix(value, `'`)) {
			value = value[1 : len(value)-1]
		}

		config.values[key] = value
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	envConfig = config
	return nil
}

// Env gets an environment variable value
func Env(key string, defaultValue ...string) string {
	// First try to get from .env file
	if envConfig != nil {
		if value, exists := envConfig.values[key]; exists {
			return value
		}
	}

	// Then try system environment variables
	if value := os.Getenv(key); value != "" {
		return value
	}

	// Return default value if provided
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// EnvInt gets an environment variable as integer
func EnvInt(key string, defaultValue ...int) int {
	value := Env(key)
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	return intValue
}

// EnvBool gets an environment variable as boolean
func EnvBool(key string, defaultValue ...bool) bool {
	value := strings.ToLower(Env(key))
	if value == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	return value == "true" || value == "1" || value == "yes" || value == "on"
}

// AutoConnect automatically connects to database using .env configuration
func AutoConnect() error {
	// Load .env file if not already loaded
	if envConfig == nil {
		if err := LoadEnv(); err != nil {
			return fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	// Get database connection type
	dbConnection := Env("DB_CONNECTION", "pgsql")

	// Build connection config from environment variables
	config := ConnectionConfig{
		Host:     Env("DB_HOST", "localhost"),
		Port:     EnvInt("DB_PORT", getDefaultPort(dbConnection)),
		Database: Env("DB_DATABASE", ""),
		Username: Env("DB_USERNAME", ""),
		Password: Env("DB_PASSWORD", ""),
		Charset:  Env("DB_CHARSET", ""),
		Options:  make(map[string]string),
	}

	// Validate required fields
	if config.Database == "" {
		return fmt.Errorf("DB_DATABASE is required in .env file or environment variables")
	}
	if config.Username == "" {
		return fmt.Errorf("DB_USERNAME is required in .env file or environment variables")
	}

	// Connect based on DB_CONNECTION type
	switch dbConnection {
	case "pgsql", "postgres", "postgresql":
		return PostgreSQL(config)
	case "mysql":
		return MySQL(config)
	default:
		return fmt.Errorf("unsupported DB_CONNECTION type: %s (supported: pgsql, mysql)", dbConnection)
	}
}

// getDefaultPort returns the default port for a database connection type
func getDefaultPort(dbConnection string) int {
	switch dbConnection {
	case "pgsql", "postgres", "postgresql":
		return 5432
	case "mysql":
		return 3306
	default:
		return 5432
	}
}

// Init initializes the database connection automatically
// This function should be called at the beginning of your application
func Init() error {
	return AutoConnect()
}
