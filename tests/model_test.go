package tests

import (
	"testing"
	"time"

	"github.com/crashana/go-eloquent"
	"github.com/crashana/go-eloquent/tests/models"
)

func setupTestDB(t *testing.T) {
	// Set up in-memory SQLite database for testing
	err := eloquent.SQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	// Create test tables
	conn := eloquent.DB()
	if conn == nil {
		t.Fatal("Failed to get database connection")
	}

	// Create users table
	_, err = conn.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			email_verified_at DATETIME,
			is_admin BOOLEAN DEFAULT FALSE,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	// Create posts table
	_, err = conn.Exec(`
		CREATE TABLE posts (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT,
			user_id TEXT,
			published BOOLEAN DEFAULT FALSE,
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
			id TEXT PRIMARY KEY,
			user_id TEXT UNIQUE,
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
}

func teardownTestDB() {
	_ = eloquent.GetManager().CloseAll()
}

func TestModelCreate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Test creating a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "password123",
		"is_admin": false,
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created with correct attributes
	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %s", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got %s", user.Email)
	}
	if user.IsAdmin != false {
		t.Errorf("Expected is_admin false, got %t", user.IsAdmin)
	}
	if user.Status != "active" {
		t.Errorf("Expected status 'active', got %s", user.Status)
	}

	// Verify ID was generated
	if user.ID == "" {
		t.Error("Expected ID to be generated, got empty string")
	}

	// Verify timestamps were set
	if user.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}
}

func TestModelFind(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user first
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Jane Doe",
		"email":    "jane@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test finding the user by ID
	foundUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	// Verify found user has correct attributes
	if foundUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, foundUser.ID)
	}
	if foundUser.Name != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got %s", foundUser.Name)
	}
	if foundUser.Email != "jane@example.com" {
		t.Errorf("Expected email 'jane@example.com', got %s", foundUser.Email)
	}
}

func TestModelFirst(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create multiple users
	_, err := models.User.Create(map[string]interface{}{
		"name":     "First User",
		"email":    "first@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	_, err = models.User.Create(map[string]interface{}{
		"name":     "Second User",
		"email":    "second@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create second user: %v", err)
	}

	// Test First method
	user, err := models.User.First()
	if err != nil {
		t.Fatalf("Failed to get first user: %v", err)
	}

	if user.Name != "First User" {
		t.Errorf("Expected first user name 'First User', got %s", user.Name)
	}
}

func TestModelAll(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create multiple users
	_, err := models.User.Create(map[string]interface{}{
		"name":     "User 1",
		"email":    "user1@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	_, err = models.User.Create(map[string]interface{}{
		"name":     "User 2",
		"email":    "user2@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	// Test All method
	users, err := models.User.All()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Verify users have correct attributes
	if users[0].Name != "User 1" {
		t.Errorf("Expected first user name 'User 1', got %s", users[0].Name)
	}
	if users[1].Name != "User 2" {
		t.Errorf("Expected second user name 'User 2', got %s", users[1].Name)
	}
}

func TestModelWhere(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create users with different statuses
	_, err := models.User.Create(map[string]interface{}{
		"name":     "Active User",
		"email":    "active@example.com",
		"password": "password123",
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create active user: %v", err)
	}

	_, err = models.User.Create(map[string]interface{}{
		"name":     "Inactive User",
		"email":    "inactive@example.com",
		"password": "password123",
		"status":   "inactive",
	})
	if err != nil {
		t.Fatalf("Failed to create inactive user: %v", err)
	}

	// Test Where method
	activeUsers, err := models.User.Where("status", "active").Get()
	if err != nil {
		t.Fatalf("Failed to get active users: %v", err)
	}

	if len(activeUsers) != 1 {
		t.Errorf("Expected 1 active user, got %d", len(activeUsers))
	}

	if activeUsers[0].Name != "Active User" {
		t.Errorf("Expected active user name 'Active User', got %s", activeUsers[0].Name)
	}

	// Test Where with First
	inactiveUser, err := models.User.Where("status", "inactive").First()
	if err != nil {
		t.Fatalf("Failed to get inactive user: %v", err)
	}

	if inactiveUser.Name != "Inactive User" {
		t.Errorf("Expected inactive user name 'Inactive User', got %s", inactiveUser.Name)
	}
}

func TestModelUpdate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Original Name",
		"email":    "original@example.com",
		"password": "password123",
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure updated_at changes
	time.Sleep(10 * time.Millisecond)

	// Test Update method
	err = user.Update(map[string]interface{}{
		"name":   "Updated Name",
		"status": "inactive",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify attributes were updated
	if user.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got %s", user.Name)
	}
	if user.Status != "inactive" {
		t.Errorf("Expected status 'inactive', got %s", user.Status)
	}

	// Verify updated_at was changed
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to be updated")
	}

	// Verify changes were persisted to database
	foundUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find updated user: %v", err)
	}

	if foundUser.Name != "Updated Name" {
		t.Errorf("Expected persisted name 'Updated Name', got %s", foundUser.Name)
	}
	if foundUser.Status != "inactive" {
		t.Errorf("Expected persisted status 'inactive', got %s", foundUser.Status)
	}
}

func TestModelSave(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Original Name",
		"email":    "original@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure updated_at changes
	time.Sleep(10 * time.Millisecond)

	// Test direct attribute modification + Save
	user.Name = "Saved Name"
	user.Status = "premium"

	err = user.Save()
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Verify updated_at was changed
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to be updated")
	}

	// Verify changes were persisted to database
	foundUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find saved user: %v", err)
	}

	if foundUser.Name != "Saved Name" {
		t.Errorf("Expected persisted name 'Saved Name', got %s", foundUser.Name)
	}
	if foundUser.Status != "premium" {
		t.Errorf("Expected persisted status 'premium', got %s", foundUser.Status)
	}
}

func TestModelDelete(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "To Be Deleted",
		"email":    "delete@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userID := user.ID

	// Test Delete method
	err = user.Delete()
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify user was deleted from database
	_, err = models.User.Find(userID)
	if err == nil {
		t.Error("Expected error when finding deleted user, got nil")
	}

	// Verify user count is 0
	users, err := models.User.All()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users after deletion, got %d", len(users))
	}
}

func TestModelFillable(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user with fillable attributes
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Fillable User",
		"email":    "fillable@example.com",
		"password": "password123",
		"is_admin": true,
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify fillable attributes were set
	if user.Name != "Fillable User" {
		t.Errorf("Expected name 'Fillable User', got %s", user.Name)
	}
	if user.Email != "fillable@example.com" {
		t.Errorf("Expected email 'fillable@example.com', got %s", user.Email)
	}
	if user.IsAdmin != true {
		t.Errorf("Expected is_admin true, got %t", user.IsAdmin)
	}
	if user.Status != "active" {
		t.Errorf("Expected status 'active', got %s", user.Status)
	}
}

func TestModelCasts(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user with boolean cast
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Cast User",
		"email":    "cast@example.com",
		"password": "password123",
		"is_admin": true,
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify boolean cast works
	if user.IsAdmin != true {
		t.Errorf("Expected is_admin true, got %t", user.IsAdmin)
	}

	// Test datetime cast
	if user.CreatedAt.IsZero() {
		t.Error("Expected created_at to be cast to time.Time")
	}
}

func TestModelRelationships(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Author User",
		"email":    "author@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a post for the user
	post, err := models.Post.Create(map[string]interface{}{
		"title":     "Test Post",
		"content":   "This is a test post",
		"user_id":   user.ID,
		"published": true,
	})
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	// Verify post was created correctly
	if post.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got %s", post.Title)
	}
	if post.UserID != user.ID {
		t.Errorf("Expected user_id %s, got %s", user.ID, post.UserID)
	}
	if post.Published != true {
		t.Errorf("Expected published true, got %t", post.Published)
	}

	// Create a profile for the user
	profile, err := models.Profile.Create(map[string]interface{}{
		"user_id": user.ID,
		"bio":     "Test bio",
		"avatar":  "avatar.jpg",
	})
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Verify profile was created correctly
	if profile.UserID != user.ID {
		t.Errorf("Expected user_id %s, got %s", user.ID, profile.UserID)
	}
	if profile.Bio != "Test bio" {
		t.Errorf("Expected bio 'Test bio', got %s", profile.Bio)
	}
}

func TestModelChainedQueries(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB()

	// Create users with different attributes
	_, err := models.User.Create(map[string]interface{}{
		"name":     "Admin User",
		"email":    "admin@example.com",
		"password": "password123",
		"is_admin": true,
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	_, err = models.User.Create(map[string]interface{}{
		"name":     "Regular User",
		"email":    "regular@example.com",
		"password": "password123",
		"is_admin": false,
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create regular user: %v", err)
	}

	_, err = models.User.Create(map[string]interface{}{
		"name":     "Inactive Admin",
		"email":    "inactive@example.com",
		"password": "password123",
		"is_admin": true,
		"status":   "inactive",
	})
	if err != nil {
		t.Fatalf("Failed to create inactive admin: %v", err)
	}

	// Test chained where clauses
	activeAdmins, err := models.User.Where("is_admin", true).Where("status", "active").Get()
	if err != nil {
		t.Fatalf("Failed to get active admins: %v", err)
	}

	if len(activeAdmins) != 1 {
		t.Errorf("Expected 1 active admin, got %d", len(activeAdmins))
	}

	if activeAdmins[0].Name != "Admin User" {
		t.Errorf("Expected admin name 'Admin User', got %s", activeAdmins[0].Name)
	}

	// Test chained where with First
	regularUser, err := models.User.Where("is_admin", false).Where("status", "active").First()
	if err != nil {
		t.Fatalf("Failed to get regular user: %v", err)
	}

	if regularUser.Name != "Regular User" {
		t.Errorf("Expected regular user name 'Regular User', got %s", regularUser.Name)
	}
}
