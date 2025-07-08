package eloquent

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Connection represents a database connection
type Connection struct {
	DB     *sqlx.DB
	Driver string
	Name   string
}

// ConnectionConfig holds database connection configuration
type ConnectionConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Charset  string
	Options  map[string]string
}

// ConnectionManager manages database connections
type ConnectionManager struct {
	connections map[string]*Connection
	default_    string
}

var manager *ConnectionManager

// NewConnectionManager creates a new connection manager
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
		default_:    "default",
	}
}

// GetManager returns the global connection manager
func GetManager() *ConnectionManager {
	if manager == nil {
		manager = NewConnectionManager()
	}
	return manager
}

// AddConnection adds a new database connection
func (cm *ConnectionManager) AddConnection(name string, config ConnectionConfig) error {
	dsn, err := buildDSN(config)
	if err != nil {
		return err
	}

	db, err := sqlx.Connect(config.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	cm.connections[name] = &Connection{
		DB:     db,
		Driver: config.Driver,
		Name:   name,
	}

	return nil
}

// GetConnection returns a database connection by name
func (cm *ConnectionManager) GetConnection(name ...string) *Connection {
	connName := cm.default_
	if len(name) > 0 && name[0] != "" {
		connName = name[0]
	}

	if conn, exists := cm.connections[connName]; exists {
		return conn
	}

	// Return default connection if exists
	if conn, exists := cm.connections[cm.default_]; exists {
		return conn
	}

	return nil
}

// SetDefaultConnection sets the default connection name
func (cm *ConnectionManager) SetDefaultConnection(name string) {
	cm.default_ = name
}

// CloseAll closes all database connections
func (cm *ConnectionManager) CloseAll() error {
	var errs []string

	for name, conn := range cm.connections {
		if err := conn.DB.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("failed to close connection '%s': %v", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %s", strings.Join(errs, "; "))
	}

	return nil
}

// Connection methods

// Select executes a select query and returns the results
func (c *Connection) Select(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return c.scanRows(rows)
}

// Insert executes an insert query
func (c *Connection) Insert(query string, args ...interface{}) (sql.Result, error) {
	return c.DB.Exec(query, args...)
}

// Update executes an update query
func (c *Connection) Update(query string, args ...interface{}) (sql.Result, error) {
	return c.DB.Exec(query, args...)
}

// Delete executes a delete query
func (c *Connection) Delete(query string, args ...interface{}) (sql.Result, error) {
	return c.DB.Exec(query, args...)
}

// Exec executes a query without returning rows
func (c *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.DB.Exec(query, args...)
}

// Begin starts a new transaction
func (c *Connection) Begin() (*sqlx.Tx, error) {
	return c.DB.Beginx()
}

// Transaction executes a function within a transaction
func (c *Connection) Transaction(fn func(*sqlx.Tx) error) error {
	tx, err := c.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}

// scanRows converts sql.Rows to []map[string]interface{}
func (c *Connection) scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

// buildDSN builds a database connection string based on the driver
func buildDSN(config ConnectionConfig) (string, error) {
	switch config.Driver {
	case "mysql":
		return buildMySQLDSN(config), nil
	case "postgres":
		return buildPostgresDSN(config), nil
	case "sqlite3":
		return buildSQLiteDSN(config), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
}

// buildMySQLDSN builds MySQL connection string
func buildMySQLDSN(config ConnectionConfig) string {
	charset := config.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		charset,
	)

	for key, value := range config.Options {
		dsn += fmt.Sprintf("&%s=%s", key, value)
	}

	return dsn
}

// buildPostgresDSN builds PostgreSQL connection string
func buildPostgresDSN(config ConnectionConfig) string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Database,
	)

	for key, value := range config.Options {
		dsn += fmt.Sprintf(" %s=%s", key, value)
	}

	return dsn
}

// buildSQLiteDSN builds SQLite connection string
func buildSQLiteDSN(config ConnectionConfig) string {
	return config.Database
}

// Quick setup functions

// MySQL creates a MySQL connection
func MySQL(config ConnectionConfig) error {
	config.Driver = "mysql"
	if config.Port == 0 {
		config.Port = 3306
	}
	return GetManager().AddConnection("default", config)
}

// PostgreSQL creates a PostgreSQL connection
func PostgreSQL(config ConnectionConfig) error {
	config.Driver = "postgres"
	if config.Port == 0 {
		config.Port = 5432
	}
	return GetManager().AddConnection("default", config)
}

// SQLite creates a SQLite connection
func SQLite(database string) error {
	config := ConnectionConfig{
		Driver:   "sqlite3",
		Database: database,
	}
	return GetManager().AddConnection("default", config)
}

// DB returns the default database connection
func DB(name ...string) *Connection {
	return GetManager().GetConnection(name...)
}

// init automatically initializes database connection from .env file
func init() {
	// Try to auto-connect from .env file
	// This is optional - if it fails, user can still manually connect
	_ = AutoConnect() // Ignore errors - user can manually connect if needed
}
