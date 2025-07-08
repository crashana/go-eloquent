package main

import (
	"fmt"
	"log"

	"github.com/crashana/go-eloquent"
)

// User model example
type User struct {
	*eloquent.BaseModel
}

// Static methods for User model
func UserQuery() *eloquent.ModelQueryBuilder {
	user := NewUser()
	return user.Query()
}

func UserWhere(column string, args ...interface{}) *eloquent.ModelQueryBuilder {
	return UserQuery().Where(column, args...)
}

func UserAll() ([]eloquent.Model, error) {
	return UserQuery().Get()
}

func UserFirst() (eloquent.Model, error) {
	return UserQuery().First()
}

func UserFind(id interface{}) (eloquent.Model, error) {
	return UserQuery().Find(id)
}

// NewUser creates a new User instance
func NewUser() *User {
	user := &User{
		BaseModel: eloquent.NewBaseModel(),
	}

	// Configure the model
	user.Table("users").
		Fillable("name", "email", "password").
		Hidden("password", "remember_token").
		Casts(map[string]string{
			"email_verified_at": "datetime",
			"created_at":        "datetime",
			"updated_at":        "datetime",
		})

	return user
}

// Define relationships
func (u *User) Posts() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(u)
	return rb.HasMany("posts", "Post")
}

func (u *User) Profile() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(u)
	return rb.HasOne("profile", "Profile")
}

// Post model example
type Post struct {
	*eloquent.BaseModel
}

func NewPost() *Post {
	post := &Post{
		BaseModel: eloquent.NewBaseModel(),
	}

	post.Table("posts").
		Fillable("title", "content", "user_id", "published_at").
		Casts(map[string]string{
			"published_at": "datetime",
			"created_at":   "datetime",
			"updated_at":   "datetime",
		})

	return post
}

func (p *Post) Author() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(p)
	return rb.BelongsTo("author", "User")
}

func (p *Post) Tags() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(p)
	return rb.BelongsToMany("tags", "Tag", "post_tag")
}

// Profile model example
type Profile struct {
	*eloquent.BaseModel
}

func NewProfile() *Profile {
	profile := &Profile{
		BaseModel: eloquent.NewBaseModel(),
	}

	profile.Table("profiles").
		Fillable("user_id", "first_name", "last_name", "bio", "avatar")

	return profile
}

func (p *Profile) User() *eloquent.Relationship {
	rb := eloquent.NewRelationshipBuilder(p)
	return rb.BelongsTo("user", "User")
}

func main() {
	// Setup database connection
	err := eloquent.SQLite("example.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer eloquent.GetManager().CloseAll()

	fmt.Println("=== Go Eloquent ORM Example ===")

	// Example 1: Creating a new user
	fmt.Println("\n1. Creating a new user:")
	user := NewUser()
	user.Fill(map[string]interface{}{
		"name":     "John Doe",
		"email":    "john@example.com",
		"password": "secret123",
	})

	fmt.Printf("User attributes: %v\n", user.ToMap())
	fmt.Printf("Is dirty: %v\n", user.IsDirty())

	// Example 2: Query Builder usage
	fmt.Println("\n2. Query Builder examples:")

	// Get connection and create query builder
	db := eloquent.DB()
	if db != nil {
		// Basic queries
		qb := eloquent.NewQueryBuilder(db)

		// Demonstrate query building (would execute with real DB)
		sql, args := qb.Table("users").
			Select("id", "name", "email").
			Where("active", true).
			Where("created_at", ">=", "2023-01-01").
			OrderBy("name", "asc").
			Limit(10).
			ToSQL()

		fmt.Printf("Generated SQL: %s\n", sql)
		fmt.Printf("Args: %v\n", args)

		// Complex query example
		sql2, args2 := qb.Table("posts").
			Select("posts.*", "users.name as author_name").
			Join("users", "posts.user_id", "=", "users.id").
			Where("posts.published", true).
			WhereIn("posts.category_id", []interface{}{1, 2, 3}).
			GroupBy("posts.user_id").
			Having("COUNT(*)", ">", 5).
			OrderByDesc("posts.created_at").
			ToSQL()

		fmt.Printf("\nComplex SQL: %s\n", sql2)
		fmt.Printf("Args: %v\n", args2)
	}

	// Example 3: Scopes usage
	fmt.Println("\n3. Scopes examples:")

	// Create some example scopes
	activeScope := eloquent.WhereStatusScope("active")
	recentScope := eloquent.RecentScope(30) // Last 30 days
	publishedScope := eloquent.PublishedScope()

	// Chain scopes
	combinedScope := eloquent.ChainScopes(activeScope, recentScope, publishedScope)

	fmt.Printf("Created combined scope with: active status + recent + published (scope count: %d)\n", 3)
	_ = combinedScope // Demonstrate scope creation

	// Example 4: Model states
	fmt.Println("\n4. Model state examples:")

	user.SetAttribute("name", "Jane Doe")
	fmt.Printf("After changing name - Is dirty: %v\n", user.IsDirty())
	fmt.Printf("Dirty attributes: %v\n", user.GetDirty())
	fmt.Printf("Original name: %v\n", user.GetOriginal("name"))
	fmt.Printf("Current name: %v\n", user.GetAttribute("name"))

	// Example 5: Casting
	fmt.Println("\n5. Attribute casting examples:")

	post := NewPost()
	post.SetAttribute("published_at", "2023-12-01 10:00:00")
	fmt.Printf("Published at (casted): %v\n", post.GetAttribute("published_at"))

	// Example 6: Relationships (structure demonstration)
	fmt.Println("\n6. Relationship definitions:")

	// Show relationship types
	userPosts := user.Posts()
	userProfile := user.Profile()
	postAuthor := NewPost().Author()
	postTags := NewPost().Tags()

	fmt.Printf("User->Posts relationship type: %s\n", userPosts.Type)
	fmt.Printf("User->Profile relationship type: %s\n", userProfile.Type)
	fmt.Printf("Post->Author relationship type: %s\n", postAuthor.Type)
	fmt.Printf("Post->Tags relationship type: %s\n", postTags.Type)

	// Example 7: Pagination demo
	fmt.Println("\n7. Pagination example:")

	if db != nil {
		qb := eloquent.NewQueryBuilder(db)
		// This would work with a real database connection
		fmt.Println("Pagination query structure:")
		sql, args := qb.Table("users").
			Where("active", true).
			OrderBy("created_at", "desc").
			Offset(0).
			Limit(15).
			ToSQL()

		fmt.Printf("Page 1 SQL: %s\n", sql)
		fmt.Printf("Args: %v\n", args)
	}

	// Example 8: Search functionality
	fmt.Println("\n8. Search scope example:")

	searchScope := eloquent.SearchScope("john doe", "name", "email", "bio")
	advancedSearchScope := eloquent.AdvancedSearchScope(map[string]interface{}{
		"name":   "john",
		"active": true,
		"age": map[string]interface{}{
			"min": 18,
			"max": 65,
		},
		"categories": []interface{}{1, 2, 3},
	})

	fmt.Println("Created search scopes for text and advanced filtering")
	_ = searchScope         // Demonstrate search scope creation
	_ = advancedSearchScope // Demonstrate advanced search scope creation

	// Example 9: Connection management
	fmt.Println("\n9. Connection management:")

	// Add a second connection (example)
	mysqlConfig := eloquent.ConnectionConfig{
		Host:     "localhost",
		Port:     3306,
		Database: "test_db",
		Username: "user",
		Password: "password",
	}

	fmt.Printf("MySQL config example: %+v\n", mysqlConfig)

	// Example 10: Soft deletes
	fmt.Println("\n10. Soft delete example:")

	user.SetAttribute("deleted_at", nil) // Not deleted
	fmt.Printf("User uses soft deletes: %v\n", user.GetDeletedAtColumn() != "")

	// Simulate soft delete
	user.Delete() // This would set deleted_at timestamp
	fmt.Printf("After soft delete, deleted_at: %v\n", user.GetAttribute("deleted_at"))

	fmt.Println("\n=== Example completed successfully! ===")
	fmt.Println("\nThis example demonstrates the Laravel Eloquent-like API in Go.")
	fmt.Println("To use with a real database, configure your connection and run migrations.")

	// Example 11: Laravel-style Model Querying
	fmt.Println("\n11. Laravel-style Model Querying:")

	// Note: These examples show the syntax, but would need a real database connection to execute
	fmt.Println("\nExample syntax for Laravel-style querying:")

	// Show how to use the new Laravel-style methods
	fmt.Println(`
	// Find user by email
	user, err := UserWhere("email", "testUser@gmail.com").First()
	
	// Find user by ID
	user, err := UserFind(1)
	
	// Get all active users
	users, err := UserWhere("active", true).Get()
	
	// Complex query
	users, err := UserWhere("age", ">", 18).
		Where("status", "active").
		OrderBy("created_at", "desc").
		Limit(10).
		Get()
	
	// Using instance methods
	userInstance := NewUser()
	user, err := userInstance.Where("email", "test@example.com").First()
	`)

	// Demonstrate the query building (without execution)
	fmt.Println("Query building examples:")

	// Show SQL generation for Laravel-style queries
	userQuery := UserWhere("email", "testUser@gmail.com")
	sql, args := userQuery.QueryBuilder.ToSQL()
	fmt.Printf("Laravel-style query SQL: %s\n", sql)
	fmt.Printf("Args: %v\n", args)

	// Complex query example
	complexQuery := UserWhere("age", ">", 18).
		Where("status", "active").
		OrderBy("created_at", "desc").
		Limit(10)
	sql2, args2 := complexQuery.QueryBuilder.ToSQL()
	fmt.Printf("Complex query SQL: %s\n", sql2)
	fmt.Printf("Args: %v\n", args2)
}
