package eloquent

import (
	"testing"
)

func setupRelationshipTestDB(t *testing.T) {
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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create posts table: %v", err)
	}

	// Create profiles table
	_, err = conn.Exec(`
		CREATE TABLE profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER UNIQUE,
			bio TEXT,
			avatar TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create profiles table: %v", err)
	}

	// Create tags table
	_, err = conn.Exec(`
		CREATE TABLE tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create tags table: %v", err)
	}

	// Create post_tags pivot table
	_, err = conn.Exec(`
		CREATE TABLE post_tags (
			post_id INTEGER,
			tag_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create post_tags table: %v", err)
	}
}

func teardownRelationshipTestDB() {
	GetManager().CloseAll()
}

func TestRelationshipBuilder(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock model for testing
	model := NewBaseModel()
	model.Table("users").PrimaryKey("id")

	// Test creating relationship builder
	rb := NewRelationshipBuilder(model)
	if rb == nil {
		t.Fatal("Expected relationship builder to be created")
	}
}

func TestBelongsToRelationship(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock post model
	postModel := NewBaseModel()
	postModel.Table("posts").PrimaryKey("id")

	// Create BelongsTo relationship
	rb := NewRelationshipBuilder(postModel)
	relationship := rb.BelongsTo("user", "users")

	if relationship == nil {
		t.Fatal("Expected BelongsTo relationship to be created")
	}

	if relationship.Type != "belongsTo" {
		t.Errorf("Expected relationship type 'belongsTo', got %s", relationship.Type)
	}

	if relationship.Related != "users" {
		t.Errorf("Expected related 'users', got %s", relationship.Related)
	}

	if relationship.ForeignKey != "users_id" {
		t.Errorf("Expected foreign key 'users_id', got %s", relationship.ForeignKey)
	}

	if relationship.LocalKey != "id" {
		t.Errorf("Expected local key 'id', got %s", relationship.LocalKey)
	}
}

func TestHasOneRelationship(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock user model
	userModel := NewBaseModel()
	userModel.Table("users").PrimaryKey("id")

	// Create HasOne relationship
	rb := NewRelationshipBuilder(userModel)
	relationship := rb.HasOne("profile", "profiles")

	if relationship == nil {
		t.Fatal("Expected HasOne relationship to be created")
	}

	if relationship.Type != "hasOne" {
		t.Errorf("Expected relationship type 'hasOne', got %s", relationship.Type)
	}

	if relationship.Related != "profiles" {
		t.Errorf("Expected related 'profiles', got %s", relationship.Related)
	}

	if relationship.ForeignKey != "users_id" {
		t.Errorf("Expected foreign key 'users_id', got %s", relationship.ForeignKey)
	}

	if relationship.LocalKey != "id" {
		t.Errorf("Expected local key 'id', got %s", relationship.LocalKey)
	}
}

func TestHasManyRelationship(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock user model
	userModel := NewBaseModel()
	userModel.Table("users").PrimaryKey("id")

	// Create HasMany relationship
	rb := NewRelationshipBuilder(userModel)
	relationship := rb.HasMany("posts", "posts")

	if relationship == nil {
		t.Fatal("Expected HasMany relationship to be created")
	}

	if relationship.Type != "hasMany" {
		t.Errorf("Expected relationship type 'hasMany', got %s", relationship.Type)
	}

	if relationship.Related != "posts" {
		t.Errorf("Expected related 'posts', got %s", relationship.Related)
	}

	if relationship.ForeignKey != "users_id" {
		t.Errorf("Expected foreign key 'users_id', got %s", relationship.ForeignKey)
	}

	if relationship.LocalKey != "id" {
		t.Errorf("Expected local key 'id', got %s", relationship.LocalKey)
	}
}

func TestBelongsToManyRelationship(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock post model
	postModel := NewBaseModel()
	postModel.Table("posts").PrimaryKey("id")

	// Create BelongsToMany relationship
	rb := NewRelationshipBuilder(postModel)
	relationship := rb.BelongsToMany("tags", "tags", "post_tags")

	if relationship == nil {
		t.Fatal("Expected BelongsToMany relationship to be created")
	}

	if relationship.Type != "belongsToMany" {
		t.Errorf("Expected relationship type 'belongsToMany', got %s", relationship.Type)
	}

	if relationship.Related != "tags" {
		t.Errorf("Expected related 'tags', got %s", relationship.Related)
	}

	if relationship.PivotTable != "post_tags" {
		t.Errorf("Expected pivot table 'post_tags', got %s", relationship.PivotTable)
	}

	if relationship.FirstKey != "posts_id" {
		t.Errorf("Expected first key 'posts_id', got %s", relationship.FirstKey)
	}

	if relationship.SecondKey != "tags_id" {
		t.Errorf("Expected second key 'tags_id', got %s", relationship.SecondKey)
	}
}

func TestRelationshipConstraints(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock model
	model := NewBaseModel()
	model.Table("users").PrimaryKey("id")

	// Test relationship with constraints
	rb := NewRelationshipBuilder(model)
	relationship := rb.HasMany("published_posts", "posts").
		Where("published", true)

	if len(relationship.Constraints) != 1 {
		t.Errorf("Expected 1 constraint, got %d", len(relationship.Constraints))
	}

	// Test chaining multiple constraints
	relationship = rb.HasMany("active_posts", "posts").
		Where("status", "published").
		Where("deleted_at", nil).
		OrderBy("created_at", "DESC")

	if len(relationship.Constraints) != 3 {
		t.Errorf("Expected 3 constraints, got %d", len(relationship.Constraints))
	}
}

func TestRelationshipMethods(t *testing.T) {
	setupRelationshipTestDB(t)
	defer teardownRelationshipTestDB()

	// Create a mock model
	model := NewBaseModel()
	model.Table("users").PrimaryKey("id")

	// Test relationship method chaining
	rb := NewRelationshipBuilder(model)
	relationship := rb.HasMany("posts", "posts")

	// Test Where method
	result := relationship.Where("published", true)
	if result != relationship {
		t.Error("Expected Where to return the same relationship instance")
	}

	// Test OrderBy method
	result = relationship.OrderBy("created_at", "DESC")
	if result != relationship {
		t.Error("Expected OrderBy to return the same relationship instance")
	}

	// Test WhereIn method
	result = relationship.WhereIn("status", []interface{}{"published", "draft"})
	if result != relationship {
		t.Error("Expected WhereIn to return the same relationship instance")
	}
}

func TestRelationshipTypes(t *testing.T) {
	// Test relationship type constants
	if HasOne != "hasOne" {
		t.Errorf("Expected HasOne constant to be 'hasOne', got %s", HasOne)
	}

	if HasMany != "hasMany" {
		t.Errorf("Expected HasMany constant to be 'hasMany', got %s", HasMany)
	}

	if BelongsTo != "belongsTo" {
		t.Errorf("Expected BelongsTo constant to be 'belongsTo', got %s", BelongsTo)
	}

	if BelongsToMany != "belongsToMany" {
		t.Errorf("Expected BelongsToMany constant to be 'belongsToMany', got %s", BelongsToMany)
	}
}
