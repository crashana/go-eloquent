package eloquent

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestConnectionManager(t *testing.T) {
	// Create a new connection manager
	cm := NewConnectionManager()

	// Test initial state
	if cm.default_ != "default" {
		t.Errorf("Expected default connection name to be 'default', got %s", cm.default_)
	}

	if len(cm.connections) != 0 {
		t.Errorf("Expected empty connections map, got %d connections", len(cm.connections))
	}
}

func TestSetDefaultConnection(t *testing.T) {
	cm := NewConnectionManager()

	// Test setting default connection
	cm.SetDefaultConnection("test_connection")

	if cm.default_ != "test_connection" {
		t.Errorf("Expected default connection to be 'test_connection', got %s", cm.default_)
	}
}

func TestGetConnectionNonExistent(t *testing.T) {
	cm := NewConnectionManager()

	// Test getting non-existent connection
	conn := cm.GetConnection("nonexistent")
	if conn != nil {
		t.Error("Expected nil for non-existent connection, got connection")
	}
}

func TestAddConnectionSQLite(t *testing.T) {
	cm := NewConnectionManager()

	// Test adding SQLite connection
	config := ConnectionConfig{
		Driver:   "sqlite3",
		Database: ":memory:",
	}

	err := cm.AddConnection("sqlite_test", config)
	if err != nil {
		t.Fatalf("Failed to add SQLite connection: %v", err)
	}

	// Test that connection was added
	conn := cm.GetConnection("sqlite_test")
	if conn == nil {
		t.Error("Expected connection to be added, got nil")
	}

	if conn.Driver != "sqlite3" {
		t.Errorf("Expected driver to be 'sqlite3', got %s", conn.Driver)
	}

	if conn.Name != "sqlite_test" {
		t.Errorf("Expected connection name to be 'sqlite_test', got %s", conn.Name)
	}
}

func TestAddConnectionInvalidDriver(t *testing.T) {
	cm := NewConnectionManager()

	// Test adding connection with invalid driver
	config := ConnectionConfig{
		Driver:   "invalid_driver",
		Database: "test.db",
	}

	err := cm.AddConnection("invalid_test", config)
	if err == nil {
		t.Error("Expected error for invalid driver, got nil")
	}
}

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   ConnectionConfig
		expected string
	}{
		{
			name: "SQLite DSN",
			config: ConnectionConfig{
				Driver:   "sqlite3",
				Database: "test.db",
			},
			expected: "test.db",
		},
		{
			name: "MySQL DSN",
			config: ConnectionConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				Username: "user",
				Password: "pass",
				Charset:  "utf8mb4",
			},
			expected: "user:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "PostgreSQL DSN",
			config: ConnectionConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "user",
				Password: "pass",
			},
			expected: "host=localhost port=5432 user=user password=pass dbname=testdb sslmode=disable",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := buildDSN(test.config)
			if err != nil {
				t.Fatalf("buildDSN failed: %v", err)
			}

			if actual != test.expected {
				t.Errorf("Expected DSN: %s, got: %s", test.expected, actual)
			}
		})
	}
}

func TestBuildDSNUnsupportedDriver(t *testing.T) {
	config := ConnectionConfig{
		Driver: "unsupported",
	}

	_, err := buildDSN(config)
	if err == nil {
		t.Error("Expected error for unsupported driver, got nil")
	}
}

func TestBuildMySQLDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   ConnectionConfig
		expected string
	}{
		{
			name: "Basic MySQL DSN",
			config: ConnectionConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				Username: "user",
				Password: "pass",
			},
			expected: "user:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local",
		},
		{
			name: "MySQL DSN with custom charset",
			config: ConnectionConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				Username: "user",
				Password: "pass",
				Charset:  "utf8",
			},
			expected: "user:pass@tcp(localhost:3306)/testdb?charset=utf8&parseTime=True&loc=Local",
		},
		{
			name: "MySQL DSN with options",
			config: ConnectionConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				Username: "user",
				Password: "pass",
				Charset:  "utf8mb4",
				Options: map[string]string{
					"parseTime": "true",
					"loc":       "Local",
				},
			},
			expected: "user:pass@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local&loc=Local&parseTime=true",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := buildMySQLDSN(test.config)
			if actual != test.expected {
				t.Errorf("Expected DSN: %s, got: %s", test.expected, actual)
			}
		})
	}
}

func TestBuildPostgresDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   ConnectionConfig
		expected string
	}{
		{
			name: "Basic PostgreSQL DSN",
			config: ConnectionConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "user",
				Password: "pass",
			},
			expected: "host=localhost port=5432 user=user password=pass dbname=testdb sslmode=disable",
		},
		{
			name: "PostgreSQL DSN with options",
			config: ConnectionConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				Username: "user",
				Password: "pass",
				Options: map[string]string{
					"sslmode":  "require",
					"timezone": "UTC",
				},
			},
			expected: "host=localhost port=5432 user=user password=pass dbname=testdb sslmode=disable sslmode=require timezone=UTC",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := buildPostgresDSN(test.config)
			if actual != test.expected {
				t.Errorf("Expected DSN: %s, got: %s", test.expected, actual)
			}
		})
	}
}

func TestBuildSQLiteDSN(t *testing.T) {
	config := ConnectionConfig{
		Database: "test.db",
	}

	actual := buildSQLiteDSN(config)
	expected := "test.db"

	if actual != expected {
		t.Errorf("Expected DSN: %s, got: %s", expected, actual)
	}
}

func TestQuickSetupFunctions(t *testing.T) {
	// Test SQLite setup
	err := SQLite(":memory:")
	if err != nil {
		t.Errorf("SQLite setup failed: %v", err)
	}

	// Test that connection was created
	conn := DB()
	if conn == nil {
		t.Error("Expected connection to be created, got nil")
	}

	if conn.Driver != "sqlite3" {
		t.Errorf("Expected driver to be 'sqlite3', got %s", conn.Driver)
	}

	// Clean up
	GetManager().CloseAll()
}

func TestGetGlobalManager(t *testing.T) {
	// Test that GetManager returns the same instance
	manager1 := GetManager()
	manager2 := GetManager()

	if manager1 != manager2 {
		t.Error("Expected GetManager to return the same instance")
	}
}

func TestDBFunction(t *testing.T) {
	// Set up a test connection
	err := SQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to set up test connection: %v", err)
	}
	defer GetManager().CloseAll()

	// Test getting default connection
	conn := DB()
	if conn == nil {
		t.Error("Expected connection, got nil")
	}

	// Test getting named connection (should return default since no named connection exists)
	conn2 := DB("nonexistent")
	if conn2 == nil {
		t.Error("Expected default connection when named connection doesn't exist, got nil")
	}
}

func TestConnectionMethods(t *testing.T) {
	// Set up SQLite connection for testing
	err := SQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to set up test connection: %v", err)
	}
	defer GetManager().CloseAll()

	conn := DB()
	if conn == nil {
		t.Fatal("Expected connection, got nil")
	}

	// Test Exec method
	_, err = conn.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	// Test Insert method
	_, err = conn.Insert("INSERT INTO test (name) VALUES (?)", "test_name")
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}

	// Test Select method
	rows, err := conn.Select("SELECT * FROM test")
	if err != nil {
		t.Errorf("Select failed: %v", err)
	}

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	// Test Update method
	_, err = conn.Update("UPDATE test SET name = ? WHERE id = ?", "updated_name", 1)
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	// Test Delete method
	_, err = conn.Delete("DELETE FROM test WHERE id = ?", 1)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

func TestConnectionTransaction(t *testing.T) {
	// Set up SQLite connection for testing
	err := SQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to set up test connection: %v", err)
	}
	defer GetManager().CloseAll()

	conn := DB()
	if conn == nil {
		t.Fatal("Expected connection, got nil")
	}

	// Create test table
	_, err = conn.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Test successful transaction
	err = conn.Transaction(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("INSERT INTO test (name) VALUES (?)", "test_name")
		return err
	})
	if err != nil {
		t.Errorf("Transaction failed: %v", err)
	}

	// Verify data was inserted
	rows, err := conn.Select("SELECT COUNT(*) as count FROM test")
	if err != nil {
		t.Errorf("Select failed: %v", err)
	}

	if len(rows) != 1 {
		t.Errorf("Expected 1 row, got %d", len(rows))
	}

	// Test transaction rollback
	err = conn.Transaction(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("INSERT INTO test (name) VALUES (?)", "rollback_test")
		if err != nil {
			return err
		}
		return fmt.Errorf("intentional error to trigger rollback")
	})
	if err == nil {
		t.Error("Expected transaction to fail and rollback")
	}

	// Verify data was rolled back (should still be 1 row)
	rows, err = conn.Select("SELECT COUNT(*) as count FROM test")
	if err != nil {
		t.Errorf("Select failed: %v", err)
	}

	if len(rows) != 1 {
		t.Errorf("Expected 1 row after rollback, got %d", len(rows))
	}
}
