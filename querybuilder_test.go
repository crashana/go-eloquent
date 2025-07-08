package eloquent

import (
	"testing"
)

func setupQueryBuilderTestDB(t *testing.T) {
	// Set up in-memory SQLite database for testing
	err := SQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Create test tables
	conn := DB()
	if conn == nil {
		t.Fatal("Failed to get database connection")
	}

	// Create users table
	_, err = conn.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INTEGER,
			is_admin BOOLEAN DEFAULT FALSE,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Create posts table
	_, err = conn.Exec(`
		CREATE TABLE posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT,
			user_id INTEGER,
			published BOOLEAN DEFAULT FALSE,
			views INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create posts table: %v", err)
	}

	// Insert test data
	_, err = conn.Exec(`
		INSERT INTO users (name, email, age, is_admin, status) VALUES 
		('John Doe', 'john@example.com', 25, true, 'active'),
		('Jane Smith', 'jane@example.com', 30, false, 'active'),
		('Bob Johnson', 'bob@example.com', 35, false, 'inactive'),
		('Alice Brown', 'alice@example.com', 28, true, 'active')
	`)
	if err != nil {
		t.Fatalf("Failed to insert test users: %v", err)
	}

	_, err = conn.Exec(`
		INSERT INTO posts (title, content, user_id, published, views) VALUES 
		('First Post', 'Content of first post', 1, true, 100),
		('Second Post', 'Content of second post', 1, false, 50),
		('Third Post', 'Content of third post', 2, true, 200),
		('Fourth Post', 'Content of fourth post', 2, true, 150)
	`)
	if err != nil {
		t.Fatalf("Failed to insert test posts: %v", err)
	}
}

func teardownQueryBuilderTestDB() {
	GetManager().CloseAll()
}

func TestQueryBuilderBasicSelect(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test basic select
	results, err := qb.Table("users").Get()
	if err != nil {
		t.Fatalf("Failed to execute basic select: %v", err)
	}

	if len(results) != 4 {
		t.Errorf("Expected 4 users, got %d", len(results))
	}

	// Verify first user
	if results[0]["name"] != "John Doe" {
		t.Errorf("Expected first user name 'John Doe', got %s", results[0]["name"])
	}
	if results[0]["email"] != "john@example.com" {
		t.Errorf("Expected first user email 'john@example.com', got %s", results[0]["email"])
	}
}

func TestQueryBuilderSelectColumns(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test select specific columns
	results, err := qb.Table("users").Select("name", "email").Get()
	if err != nil {
		t.Fatalf("Failed to execute select with columns: %v", err)
	}

	if len(results) != 4 {
		t.Errorf("Expected 4 users, got %d", len(results))
	}

	// Verify only selected columns are present
	firstUser := results[0]
	if _, exists := firstUser["name"]; !exists {
		t.Error("Expected 'name' column to be present")
	}
	if _, exists := firstUser["email"]; !exists {
		t.Error("Expected 'email' column to be present")
	}
	if _, exists := firstUser["age"]; exists {
		t.Error("Expected 'age' column to NOT be present")
	}
}

func TestQueryBuilderWhere(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test where clause
	results, err := qb.Table("users").Where("status", "active").Get()
	if err != nil {
		t.Fatalf("Failed to execute where query: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 active users, got %d", len(results))
	}

	// Verify all results have active status
	for _, user := range results {
		if user["status"] != "active" {
			t.Errorf("Expected status 'active', got %s", user["status"])
		}
	}
}

func TestQueryBuilderWhereOperators(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()

	tests := []struct {
		name     string
		operator string
		value    interface{}
		expected int
	}{
		{"equals", "=", 25, 1},
		{"greater than", ">", 28, 2},
		{"less than", "<", 30, 2},
		{"greater than or equal", ">=", 30, 2},
		{"less than or equal", "<=", 28, 2},
		{"not equal", "!=", 25, 3},
		{"not equal (alternative)", "<>", 25, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			qb := NewQueryBuilder(db)
			results, err := qb.Table("users").Where("age", test.operator, test.value).Get()
			if err != nil {
				t.Fatalf("Failed to execute where query with %s: %v", test.operator, err)
			}

			if len(results) != test.expected {
				t.Errorf("Expected %d results for %s %v, got %d", test.expected, test.operator, test.value, len(results))
			}
		})
	}
}

func TestQueryBuilderWhereIn(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test WhereIn
	results, err := qb.Table("users").WhereIn("age", []interface{}{25, 30}).Get()
	if err != nil {
		t.Fatalf("Failed to execute WhereIn query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users with age 25 or 30, got %d", len(results))
	}

	// Verify results
	ages := make([]interface{}, len(results))
	for i, user := range results {
		ages[i] = user["age"]
	}

	expectedAges := []interface{}{int64(25), int64(30)}
	for _, expectedAge := range expectedAges {
		found := false
		for _, age := range ages {
			if age == expectedAge {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected age %v not found in results", expectedAge)
		}
	}
}

func TestQueryBuilderWhereNotIn(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test WhereNotIn
	results, err := qb.Table("users").WhereNotIn("age", []interface{}{25, 30}).Get()
	if err != nil {
		t.Fatalf("Failed to execute WhereNotIn query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users with age not 25 or 30, got %d", len(results))
	}

	// Verify results don't contain excluded ages
	for _, user := range results {
		age := user["age"]
		if age == int64(25) || age == int64(30) {
			t.Errorf("Found excluded age %v in results", age)
		}
	}
}

func TestQueryBuilderWhereNull(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	conn := DB()

	// Insert user with NULL age
	_, err := conn.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Null User", "null@example.com", nil)
	if err != nil {
		t.Fatalf("Failed to insert user with null age: %v", err)
	}

	qb := NewQueryBuilder(db)

	// Test WhereNull
	results, err := qb.Table("users").WhereNull("age").Get()
	if err != nil {
		t.Fatalf("Failed to execute WhereNull query: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 user with null age, got %d", len(results))
	}

	if results[0]["name"] != "Null User" {
		t.Errorf("Expected user name 'Null User', got %s", results[0]["name"])
	}
}

func TestQueryBuilderWhereNotNull(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	conn := DB()

	// Insert user with NULL age
	_, err := conn.Exec("INSERT INTO users (name, email, age) VALUES (?, ?, ?)", "Null User", "null@example.com", nil)
	if err != nil {
		t.Fatalf("Failed to insert user with null age: %v", err)
	}

	qb := NewQueryBuilder(db)

	// Test WhereNotNull
	results, err := qb.Table("users").WhereNotNull("age").Get()
	if err != nil {
		t.Fatalf("Failed to execute WhereNotNull query: %v", err)
	}

	if len(results) != 4 {
		t.Errorf("Expected 4 users with non-null age, got %d", len(results))
	}

	// Verify all results have non-null age
	for _, user := range results {
		if user["age"] == nil {
			t.Error("Found user with null age in WhereNotNull results")
		}
	}
}

func TestQueryBuilderWhereBetween(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test WhereBetween
	results, err := qb.Table("users").WhereBetween("age", 25, 30).Get()
	if err != nil {
		t.Fatalf("Failed to execute WhereBetween query: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 users with age between 25 and 30, got %d", len(results))
	}

	// Verify all results are within range
	for _, user := range results {
		age := user["age"].(int64)
		if age < 25 || age > 30 {
			t.Errorf("Found user with age %d outside range 25-30", age)
		}
	}
}

func TestQueryBuilderOrWhere(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test OrWhere
	results, err := qb.Table("users").Where("age", 25).OrWhere("age", 35).Get()
	if err != nil {
		t.Fatalf("Failed to execute OrWhere query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users with age 25 or 35, got %d", len(results))
	}

	// Verify results
	ages := make([]int64, len(results))
	for i, user := range results {
		ages[i] = user["age"].(int64)
	}

	expectedAges := []int64{25, 35}
	for _, expectedAge := range expectedAges {
		found := false
		for _, age := range ages {
			if age == expectedAge {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected age %d not found in results", expectedAge)
		}
	}
}

func TestQueryBuilderOrderBy(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test OrderBy ASC
	results, err := qb.Table("users").OrderBy("age", "ASC").Get()
	if err != nil {
		t.Fatalf("Failed to execute OrderBy ASC query: %v", err)
	}

	if len(results) != 4 {
		t.Errorf("Expected 4 users, got %d", len(results))
	}

	// Verify ascending order
	for i := 1; i < len(results); i++ {
		prevAge := results[i-1]["age"].(int64)
		currAge := results[i]["age"].(int64)
		if prevAge > currAge {
			t.Errorf("Results not in ascending order: %d > %d", prevAge, currAge)
		}
	}

	// Test OrderBy DESC
	qb = NewQueryBuilder(db)
	results, err = qb.Table("users").OrderBy("age", "DESC").Get()
	if err != nil {
		t.Fatalf("Failed to execute OrderBy DESC query: %v", err)
	}

	// Verify descending order
	for i := 1; i < len(results); i++ {
		prevAge := results[i-1]["age"].(int64)
		currAge := results[i]["age"].(int64)
		if prevAge < currAge {
			t.Errorf("Results not in descending order: %d < %d", prevAge, currAge)
		}
	}
}

func TestQueryBuilderOrderByDesc(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test OrderByDesc
	results, err := qb.Table("users").OrderByDesc("age").Get()
	if err != nil {
		t.Fatalf("Failed to execute OrderByDesc query: %v", err)
	}

	// Verify descending order
	for i := 1; i < len(results); i++ {
		prevAge := results[i-1]["age"].(int64)
		currAge := results[i]["age"].(int64)
		if prevAge < currAge {
			t.Errorf("Results not in descending order: %d < %d", prevAge, currAge)
		}
	}
}

func TestQueryBuilderLimit(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Limit
	results, err := qb.Table("users").Limit(2).Get()
	if err != nil {
		t.Fatalf("Failed to execute Limit query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users with limit, got %d", len(results))
	}
}

func TestQueryBuilderOffset(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Offset
	results, err := qb.Table("users").OrderBy("id", "ASC").Offset(2).Get()
	if err != nil {
		t.Fatalf("Failed to execute Offset query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users with offset, got %d", len(results))
	}

	// Verify we got the last 2 users
	firstResult := results[0]
	if firstResult["name"] != "Bob Johnson" {
		t.Errorf("Expected first result name 'Bob Johnson', got %s", firstResult["name"])
	}
}

func TestQueryBuilderLimitOffset(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Limit with Offset
	results, err := qb.Table("users").OrderBy("id", "ASC").Limit(2).Offset(1).Get()
	if err != nil {
		t.Fatalf("Failed to execute Limit with Offset query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 users with limit and offset, got %d", len(results))
	}

	// Verify we got users 2 and 3
	if results[0]["name"] != "Jane Smith" {
		t.Errorf("Expected first result name 'Jane Smith', got %s", results[0]["name"])
	}
	if results[1]["name"] != "Bob Johnson" {
		t.Errorf("Expected second result name 'Bob Johnson', got %s", results[1]["name"])
	}
}

func TestQueryBuilderJoin(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Join
	results, err := qb.Table("posts").
		Select("posts.title", "users.name as author_name").
		Join("users", "posts.user_id", "=", "users.id").
		Get()
	if err != nil {
		t.Fatalf("Failed to execute Join query: %v", err)
	}

	if len(results) != 4 {
		t.Errorf("Expected 4 posts with authors, got %d", len(results))
	}

	// Verify join worked
	for _, result := range results {
		if result["title"] == nil {
			t.Error("Expected title to be present in join results")
		}
		if result["author_name"] == nil {
			t.Error("Expected author_name to be present in join results")
		}
	}

	// Verify specific result
	firstPost := results[0]
	if firstPost["title"] != "First Post" {
		t.Errorf("Expected first post title 'First Post', got %s", firstPost["title"])
	}
	if firstPost["author_name"] != "John Doe" {
		t.Errorf("Expected first post author 'John Doe', got %s", firstPost["author_name"])
	}
}

func TestQueryBuilderFirst(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test First
	result, err := qb.Table("users").OrderBy("id", "ASC").First()
	if err != nil {
		t.Fatalf("Failed to execute First query: %v", err)
	}

	if result["name"] != "John Doe" {
		t.Errorf("Expected first user name 'John Doe', got %s", result["name"])
	}
}

func TestQueryBuilderFind(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Find
	result, err := qb.Table("users").Find(1)
	if err != nil {
		t.Fatalf("Failed to execute Find query: %v", err)
	}

	if result["name"] != "John Doe" {
		t.Errorf("Expected user name 'John Doe', got %s", result["name"])
	}
	if result["id"] != int64(1) {
		t.Errorf("Expected user id 1, got %v", result["id"])
	}
}

func TestQueryBuilderCount(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Count
	count, err := qb.Table("users").Count()
	if err != nil {
		t.Fatalf("Failed to execute Count query: %v", err)
	}

	if count != 4 {
		t.Errorf("Expected count 4, got %d", count)
	}

	// Test Count with where
	qb = NewQueryBuilder(db)
	count, err = qb.Table("users").Where("status", "active").Count()
	if err != nil {
		t.Fatalf("Failed to execute Count with where query: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected active count 3, got %d", count)
	}
}

func TestQueryBuilderExists(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Exists - should return true
	exists, err := qb.Table("users").Where("name", "John Doe").Exists()
	if err != nil {
		t.Fatalf("Failed to execute Exists query: %v", err)
	}

	if !exists {
		t.Error("Expected exists to be true for John Doe")
	}

	// Test Exists - should return false
	qb = NewQueryBuilder(db)
	exists, err = qb.Table("users").Where("name", "Nonexistent User").Exists()
	if err != nil {
		t.Fatalf("Failed to execute Exists query: %v", err)
	}

	if exists {
		t.Error("Expected exists to be false for nonexistent user")
	}
}

func TestQueryBuilderChaining(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test complex chaining
	results, err := qb.Table("users").
		Select("name", "email", "age").
		Where("status", "active").
		Where("age", ">", 25).
		OrderBy("age", "ASC").
		Limit(2).
		Get()
	if err != nil {
		t.Fatalf("Failed to execute chained query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results from chained query, got %d", len(results))
	}

	// Verify results are ordered by age
	if results[0]["age"].(int64) > results[1]["age"].(int64) {
		t.Error("Results not ordered by age ascending")
	}

	// Verify all results meet criteria
	for _, result := range results {
		age := result["age"].(int64)
		if age <= 25 {
			t.Errorf("Found user with age %d, expected > 25", age)
		}
	}
}

func TestQueryBuilderDistinct(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	conn := DB()

	// Insert duplicate statuses
	_, err := conn.Exec("INSERT INTO users (name, email, status) VALUES (?, ?, ?)", "Test User", "test@example.com", "active")
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	qb := NewQueryBuilder(db)

	// Test Distinct
	results, err := qb.Table("users").Select("status").Distinct().Get()
	if err != nil {
		t.Fatalf("Failed to execute Distinct query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 distinct statuses, got %d", len(results))
	}

	// Verify distinct values
	statuses := make(map[string]bool)
	for _, result := range results {
		status := result["status"].(string)
		if statuses[status] {
			t.Errorf("Found duplicate status: %s", status)
		}
		statuses[status] = true
	}
}

func TestQueryBuilderGroupBy(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test GroupBy with Count
	results, err := qb.Table("users").
		Select("status", "COUNT(*) as count").
		GroupBy("status").
		Get()
	if err != nil {
		t.Fatalf("Failed to execute GroupBy query: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 status groups, got %d", len(results))
	}

	// Verify group counts
	for _, result := range results {
		status := result["status"].(string)
		count := result["count"].(int64)

		if status == "active" && count != 3 {
			t.Errorf("Expected 3 active users, got %d", count)
		}
		if status == "inactive" && count != 1 {
			t.Errorf("Expected 1 inactive user, got %d", count)
		}
	}
}

func TestQueryBuilderHaving(t *testing.T) {
	setupQueryBuilderTestDB(t)
	defer teardownQueryBuilderTestDB()

	db := DB()
	qb := NewQueryBuilder(db)

	// Test Having
	results, err := qb.Table("users").
		Select("status", "COUNT(*) as count").
		GroupBy("status").
		Having("COUNT(*)", ">", 1).
		Get()
	if err != nil {
		t.Fatalf("Failed to execute Having query: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 group with count > 1, got %d", len(results))
	}

	// Verify result
	result := results[0]
	if result["status"] != "active" {
		t.Errorf("Expected status 'active', got %s", result["status"])
	}
	if result["count"].(int64) != 3 {
		t.Errorf("Expected count 3, got %d", result["count"])
	}
}
