package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/crashana/go-eloquent"
	"github.com/crashana/go-eloquent/tests/models"
)

func setupIntegrationDB(t *testing.T) {
	// Set up in-memory SQLite database for integration testing
	err := eloquent.SQLite(":memory:")
	if err != nil {
		t.Fatalf("Failed to set up integration test database: %v", err)
	}

	// Create test tables
	conn := eloquent.DB()
	if conn == nil {
		t.Fatal("Failed to get database connection")
	}

	// Create users table with proper SQLite syntax
	_, err = conn.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			email_verified_at DATETIME,
			is_admin BOOLEAN DEFAULT 0,
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
			published BOOLEAN DEFAULT 0,
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

func teardownIntegrationDB() {
	_ = eloquent.GetManager().CloseAll()
}

func TestIntegrationFullCRUDWorkflow(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Test Create
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Integration Test User",
		"email":    "integration@example.com",
		"password": "password123",
		"is_admin": true,
		"status":   "active",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify user was created
	if user.Name != "Integration Test User" {
		t.Errorf("Expected name 'Integration Test User', got %s", user.Name)
	}
	if user.Email != "integration@example.com" {
		t.Errorf("Expected email 'integration@example.com', got %s", user.Email)
	}
	if user.IsAdmin != true {
		t.Errorf("Expected is_admin true, got %t", user.IsAdmin)
	}
	if user.ID == "" {
		t.Error("Expected ID to be generated")
	}

	// Test Read - Find by ID
	foundUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if foundUser.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, foundUser.ID)
	}
	if foundUser.Name != "Integration Test User" {
		t.Errorf("Expected name 'Integration Test User', got %s", foundUser.Name)
	}

	// Test Read - Query with Where
	activeUsers, err := models.User.Where("status", "active").Get()
	if err != nil {
		t.Fatalf("Failed to query active users: %v", err)
	}

	if len(activeUsers) != 1 {
		t.Errorf("Expected 1 active user, got %d", len(activeUsers))
	}

	// Test Update - Using Update method
	err = user.Update(map[string]interface{}{
		"name":   "Updated Integration User",
		"status": "premium",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify update
	if user.Name != "Updated Integration User" {
		t.Errorf("Expected updated name 'Updated Integration User', got %s", user.Name)
	}
	if user.Status != "premium" {
		t.Errorf("Expected updated status 'premium', got %s", user.Status)
	}

	// Test Update - Using Save method
	user.Name = "Saved Integration User"
	err = user.Save()
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Verify save
	savedUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find saved user: %v", err)
	}

	if savedUser.Name != "Saved Integration User" {
		t.Errorf("Expected saved name 'Saved Integration User', got %s", savedUser.Name)
	}

	// Test Delete
	err = user.Delete()
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify deletion
	_, err = models.User.Find(user.ID)
	if err == nil {
		t.Error("Expected error when finding deleted user")
	}

	// Verify no users exist
	allUsers, err := models.User.All()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(allUsers) != 0 {
		t.Errorf("Expected 0 users after deletion, got %d", len(allUsers))
	}
}

func TestIntegrationRelationshipWorkflow(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Create a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Author User",
		"email":    "author@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create posts for the user
	post1, err := models.Post.Create(map[string]interface{}{
		"title":     "First Post",
		"content":   "Content of first post",
		"user_id":   user.ID,
		"published": true,
	})
	if err != nil {
		t.Fatalf("Failed to create first post: %v", err)
	}

	post2, err := models.Post.Create(map[string]interface{}{
		"title":     "Second Post",
		"content":   "Content of second post",
		"user_id":   user.ID,
		"published": false,
	})
	if err != nil {
		t.Fatalf("Failed to create second post: %v", err)
	}

	// Create a profile for the user
	profile, err := models.Profile.Create(map[string]interface{}{
		"user_id": user.ID,
		"bio":     "Author bio",
		"avatar":  "author.jpg",
	})
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Verify relationships were created correctly
	if post1.UserID != user.ID {
		t.Errorf("Expected post1 user_id %s, got %s", user.ID, post1.UserID)
	}
	if post2.UserID != user.ID {
		t.Errorf("Expected post2 user_id %s, got %s", user.ID, post2.UserID)
	}
	if profile.UserID != user.ID {
		t.Errorf("Expected profile user_id %s, got %s", user.ID, profile.UserID)
	}

	// Test querying related data
	userPosts, err := models.Post.Where("user_id", user.ID).Get()
	if err != nil {
		t.Fatalf("Failed to get user posts: %v", err)
	}

	if len(userPosts) != 2 {
		t.Errorf("Expected 2 posts for user, got %d", len(userPosts))
	}

	// Test querying with conditions
	publishedPosts, err := models.Post.Where("user_id", user.ID).Where("published", true).Get()
	if err != nil {
		t.Fatalf("Failed to get published posts: %v", err)
	}

	if len(publishedPosts) != 1 {
		t.Errorf("Expected 1 published post, got %d", len(publishedPosts))
	}

	if publishedPosts[0].Title != "First Post" {
		t.Errorf("Expected published post title 'First Post', got %s", publishedPosts[0].Title)
	}

	// Test profile relationship
	userProfile, err := models.Profile.Where("user_id", user.ID).First()
	if err != nil {
		t.Fatalf("Failed to get user profile: %v", err)
	}

	if userProfile.Bio != "Author bio" {
		t.Errorf("Expected profile bio 'Author bio', got %s", userProfile.Bio)
	}
}

func TestIntegrationComplexQueries(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Create multiple users with different attributes
	users := []map[string]interface{}{
		{
			"name":     "Admin User",
			"email":    "admin@example.com",
			"password": "password123",
			"is_admin": true,
			"status":   "active",
		},
		{
			"name":     "Regular User 1",
			"email":    "user1@example.com",
			"password": "password123",
			"is_admin": false,
			"status":   "active",
		},
		{
			"name":     "Regular User 2",
			"email":    "user2@example.com",
			"password": "password123",
			"is_admin": false,
			"status":   "inactive",
		},
		{
			"name":     "Premium User",
			"email":    "premium@example.com",
			"password": "password123",
			"is_admin": false,
			"status":   "premium",
		},
	}

	createdUsers := make([]*models.UserModel, 0, len(users))
	for _, userData := range users {
		user, err := models.User.Create(userData)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		createdUsers = append(createdUsers, user)
	}

	// Verify users were created
	if len(createdUsers) != len(users) {
		t.Errorf("Expected %d created users, got %d", len(users), len(createdUsers))
	}

	// Test complex queries
	// 1. Find all active users
	activeUsers, err := models.User.Where("status", "active").Get()
	if err != nil {
		t.Fatalf("Failed to get active users: %v", err)
	}

	if len(activeUsers) != 2 {
		t.Errorf("Expected 2 active users, got %d", len(activeUsers))
	}

	// 2. Find all admin users
	adminUsers, err := models.User.Where("is_admin", true).Get()
	if err != nil {
		t.Fatalf("Failed to get admin users: %v", err)
	}

	if len(adminUsers) != 1 {
		t.Errorf("Expected 1 admin user, got %d", len(adminUsers))
	}

	if adminUsers[0].Name != "Admin User" {
		t.Errorf("Expected admin user name 'Admin User', got %s", adminUsers[0].Name)
	}

	// 3. Find active non-admin users
	activeRegularUsers, err := models.User.Where("status", "active").Where("is_admin", false).Get()
	if err != nil {
		t.Fatalf("Failed to get active regular users: %v", err)
	}

	if len(activeRegularUsers) != 1 {
		t.Errorf("Expected 1 active regular user, got %d", len(activeRegularUsers))
	}

	if activeRegularUsers[0].Name != "Regular User 1" {
		t.Errorf("Expected active regular user name 'Regular User 1', got %s", activeRegularUsers[0].Name)
	}

	// 4. Count total users
	totalUsers, err := models.User.All()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(totalUsers) != 4 {
		t.Errorf("Expected 4 total users, got %d", len(totalUsers))
	}

	// 5. Find first user (ordered by creation)
	firstUser, err := models.User.First()
	if err != nil {
		t.Fatalf("Failed to get first user: %v", err)
	}

	if firstUser.Name != "Admin User" {
		t.Errorf("Expected first user name 'Admin User', got %s", firstUser.Name)
	}
}

func TestIntegrationBatchOperations(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Create multiple users
	userCount := 10
	createdUsers := make([]*models.UserModel, 0, userCount)

	for i := 0; i < userCount; i++ {
		user, err := models.User.Create(map[string]interface{}{
			"name":     fmt.Sprintf("User %d", i+1),
			"email":    fmt.Sprintf("user%d@example.com", i+1),
			"password": "password123",
			"is_admin": i%3 == 0, // Every 3rd user is admin
			"status":   "active",
		})
		if err != nil {
			t.Fatalf("Failed to create user %d: %v", i+1, err)
		}
		createdUsers = append(createdUsers, user)
	}

	// Verify all users were created
	allUsers, err := models.User.All()
	if err != nil {
		t.Fatalf("Failed to get all users: %v", err)
	}

	if len(allUsers) != userCount {
		t.Errorf("Expected %d users, got %d", userCount, len(allUsers))
	}

	// Test batch update (update status for all users)
	for _, user := range createdUsers {
		err = user.Update(map[string]interface{}{
			"status": "updated",
		})
		if err != nil {
			t.Fatalf("Failed to update user %s: %v", user.ID, err)
		}
	}

	// Verify all users were updated
	updatedUsers, err := models.User.Where("status", "updated").Get()
	if err != nil {
		t.Fatalf("Failed to get updated users: %v", err)
	}

	if len(updatedUsers) != userCount {
		t.Errorf("Expected %d updated users, got %d", userCount, len(updatedUsers))
	}

	// Test batch delete (delete half of the users)
	for i := 0; i < userCount/2; i++ {
		err = createdUsers[i].Delete()
		if err != nil {
			t.Fatalf("Failed to delete user %s: %v", createdUsers[i].ID, err)
		}
	}

	// Verify remaining users
	remainingUsers, err := models.User.All()
	if err != nil {
		t.Fatalf("Failed to get remaining users: %v", err)
	}

	expectedRemaining := userCount - userCount/2
	if len(remainingUsers) != expectedRemaining {
		t.Errorf("Expected %d remaining users, got %d", expectedRemaining, len(remainingUsers))
	}
}

func TestIntegrationTransactionSimulation(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Simulate a transaction-like operation
	// Create user and profile together
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Transaction User",
		"email":    "transaction@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	_, err = models.Profile.Create(map[string]interface{}{
		"user_id": user.ID,
		"bio":     "Transaction user bio",
		"avatar":  "transaction.jpg",
	})
	if err != nil {
		// If profile creation fails, clean up user
		_ = user.Delete()
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Verify both were created successfully
	foundUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find created user: %v", err)
	}

	foundProfile, err := models.Profile.Where("user_id", user.ID).First()
	if err != nil {
		t.Fatalf("Failed to find created profile: %v", err)
	}

	if foundUser.Name != "Transaction User" {
		t.Errorf("Expected user name 'Transaction User', got %s", foundUser.Name)
	}

	if foundProfile.Bio != "Transaction user bio" {
		t.Errorf("Expected profile bio 'Transaction user bio', got %s", foundProfile.Bio)
	}

	// Test cleanup (delete both)
	err = foundProfile.Delete()
	if err != nil {
		t.Fatalf("Failed to delete profile: %v", err)
	}

	err = foundUser.Delete()
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Verify cleanup
	_, err = models.User.Find(user.ID)
	if err == nil {
		t.Error("Expected error when finding deleted user")
	}

	_, err = models.Profile.Where("user_id", user.ID).First()
	if err == nil {
		t.Error("Expected error when finding deleted profile")
	}
}

func TestIntegrationTimestamps(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Create a user
	user, err := models.User.Create(map[string]interface{}{
		"name":     "Timestamp User",
		"email":    "timestamp@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Verify timestamps were set
	if user.CreatedAt.IsZero() {
		t.Error("Expected created_at to be set")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected updated_at to be set")
	}

	// Store original timestamps
	originalCreatedAt := user.CreatedAt
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update the user
	err = user.Update(map[string]interface{}{
		"name": "Updated Timestamp User",
	})
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Verify created_at didn't change but updated_at did
	if !user.CreatedAt.Equal(originalCreatedAt) {
		t.Error("Expected created_at to remain unchanged")
	}

	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected updated_at to be updated")
	}

	// Verify timestamps are persisted
	foundUser, err := models.User.Find(user.ID)
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	if !foundUser.CreatedAt.Equal(originalCreatedAt) {
		t.Error("Expected persisted created_at to match original")
	}

	if !foundUser.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Expected persisted updated_at to be updated")
	}
}

func TestIntegrationErrorHandling(t *testing.T) {
	setupIntegrationDB(t)
	defer teardownIntegrationDB()

	// Test duplicate email constraint
	user1, err := models.User.Create(map[string]interface{}{
		"name":     "User 1",
		"email":    "duplicate@example.com",
		"password": "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	// Try to create user with duplicate email
	_, err = models.User.Create(map[string]interface{}{
		"name":     "User 2",
		"email":    "duplicate@example.com",
		"password": "password123",
	})
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}

	// Test finding non-existent user
	_, err = models.User.Find("nonexistent-id")
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}

	// Test updating non-existent user
	user1.ID = "nonexistent-id"
	err = user1.Update(map[string]interface{}{
		"name": "Updated Name",
	})
	if err == nil {
		t.Error("Expected error for updating non-existent user, got nil")
	}

	// Test deleting non-existent user
	err = user1.Delete()
	if err == nil {
		t.Error("Expected error for deleting non-existent user, got nil")
	}
}
