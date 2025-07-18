# Go Eloquent ORM

An Eloquent-inspired ORM for Go that provides the same elegant, expressive syntax for database operations.

## Features

- 🔥 **Eloquent-like API** - Familiar syntax for Laravel developers
- 🗄️ **Multiple Database Support** - MySQL, PostgreSQL, SQLite
- ⚡ **Automatic .env Configuration** - Laravel-style database setup
- 🎯 **Typed Models** - Direct attribute access without type assertions
- ✨ **Complete CRUD Operations** - Create, Read, Update, Delete with multiple methods
- 🔍 **Query Builder** - Fluent, expressive query building with method chaining
- 🔗 **Relationships** - HasOne, HasMany, BelongsTo, BelongsToMany, and more
- 🎯 **Scopes** - Reusable query constraints
- 🔄 **Soft Deletes** - Built-in soft delete functionality
- 📝 **Attribute Casting** - Automatic type conversion
- 🚀 **Mass Assignment** - Fillable and guarded attributes
- 🔑 **Auto UUID Generation** - Automatic UUID generation for PostgreSQL
- 📊 **Pagination** - Built-in pagination support
- 🔐 **Transactions** - Database transaction support

## Installation

```bash
go get github.com/crashana/go-eloquent
```

## Why Go Eloquent?

### 🚀 **Laravel-like Experience in Go**

Go Eloquent brings the beloved Laravel Eloquent ORM experience to Go, with some Go-specific improvements:

| Feature | Traditional Go ORMs | Go Eloquent |
|---------|-------------------|-------------|
| **Database Config** | Manual setup required | Automatic `.env` configuration |
| **Model Access** | Type assertions everywhere | Direct typed attribute access |
| **Query Building** | Verbose, non-chainable | Fluent, chainable methods |
| **Relationships** | Manual joins and queries | Eloquent-style relationships |
| **Attribute Casting** | Manual type conversion | Automatic casting |

### 💡 **Key Advantages**

```go
// ❌ Traditional Go ORM
type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
}

db.Select(&users, "SELECT * FROM users WHERE active = $1", true)
if user, ok := result.(*User); ok {
    fmt.Println(user.Name) // Type assertion needed
}

// ✅ Go Eloquent
user, err := models.User.Where("active", true).First()
fmt.Println(user.Name) // Direct access - no type assertions!
```

## Database Configuration

Go Eloquent supports automatic database connection from `.env` file (Laravel style) or manual configuration.

### Supported Database Connection Types

| DB_CONNECTION | Description | Default Port |
|---------------|-------------|--------------|
| `pgsql` | PostgreSQL | 5432 |
| `postgres` | PostgreSQL (alias) | 5432 |
| `postgresql` | PostgreSQL (alias) | 5432 |
| `mysql` | MySQL | 3306 |

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DB_CONNECTION` | Database type (pgsql, mysql) | No | `pgsql` |
| `DB_HOST` | Database host | No | `localhost` |
| `DB_PORT` | Database port | No | `5432` (pgsql) or `3306` (mysql) |
| `DB_DATABASE` | Database name | Yes | - |
| `DB_USERNAME` | Database username | Yes | - |
| `DB_PASSWORD` | Database password | No | - |
| `DB_CHARSET` | Database charset (MySQL only) | No | `utf8mb4` |

## Laravel vs Go Eloquent

### 🎯 **Side-by-Side Comparison**

| **Laravel (PHP)** | **Go Eloquent** |
|------------------|-----------------|
| `User::where('active', true)->first()` | `models.User.Where("active", true).First()` |
| `$user->name` | `user.Name` |
| `User::create(['name' => 'John'])` | `models.User.Create(map[string]interface{}{"name": "John"})` |
| `$user->posts()->where('published', true)->get()` | `user.Posts().Where("published", true).Get()` |

### 📝 **Configuration Comparison**

**Laravel (.env)**
```env
DB_CONNECTION=pgsql
DB_HOST=localhost
DB_DATABASE=myapp
DB_USERNAME=postgres
DB_PASSWORD=secret
```

**Go Eloquent (.env)**
```env
DB_CONNECTION=pgsql
DB_HOST=localhost
DB_DATABASE=myapp
DB_USERNAME=postgres
DB_PASSWORD=secret
```

**Laravel (PHP Code)**
```php
// No configuration needed - automatic from .env
$user = User::where('email', 'john@example.com')->first();
echo $user->name; // Direct access
```

**Go Eloquent (Go Code)**
```go
// No configuration needed - automatic from .env
user, err := models.User.Where("email", "john@example.com").First()
fmt.Println(user.Name) // Direct access - no type assertions!
```

## Quick Start

### 1. Setup Database Connection

**Option 1: Using .env file (Recommended - Laravel style)**

Create a `.env` file in your project root:

```env
# PostgreSQL Configuration
DB_CONNECTION=pgsql
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=your_database_name
DB_USERNAME=postgres
DB_PASSWORD=your_password

# MySQL Configuration (alternative)
# DB_CONNECTION=mysql
# DB_HOST=localhost
# DB_PORT=3306
# DB_DATABASE=your_database_name
# DB_USERNAME=root
# DB_PASSWORD=your_password
```

Your Go code becomes super simple:

```go
package main

import (
    "github.com/crashana/go-eloquent"
)

func main() {
    // Database connection is automatically initialized from .env file
    // No manual configuration needed!
    defer eloquent.GetManager().CloseAll()
    
    // Your models and queries work immediately
    user, err := models.User.Where("email", "john@example.com").First()
    if err != nil {
        log.Fatal(err)
    }
    
    // Direct access to typed model attributes - just like Laravel!
    fmt.Println("Welcome,", user.Name)
}
```

**Option 2: Manual Configuration (if you don't want to use .env)**

```go
package main

import (
    "github.com/crashana/go-eloquent"
)

func main() {
    // SQLite
    err := eloquent.SQLite("database.db")
    
    // MySQL
    err := eloquent.MySQL(eloquent.ConnectionConfig{
        Host:     "localhost",
        Port:     3306,
        Database: "myapp",
        Username: "user",
        Password: "password",
    })
    
    // PostgreSQL
    err := eloquent.PostgreSQL(eloquent.ConnectionConfig{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        Username: "user",
        Password: "password",
    })
    
    if err != nil {
        panic(err)
    }
    
    defer eloquent.GetManager().CloseAll()
}
```

### 2. Define Models

```go
package models

import (
    "time"
    "github.com/crashana/go-eloquent"
)

// UserModel - Typed model with direct attribute access
type UserModel struct {
    *eloquent.BaseModel
    
    // Struct fields for direct access - like Laravel Eloquent
    ID              string    `json:"id" db:"id"`
    Name            string    `json:"name" db:"name"`
    Email           string    `json:"email" db:"email"`
    EmailVerifiedAt time.Time `json:"email_verified_at" db:"email_verified_at"`
    IsAdmin         bool      `json:"is_admin" db:"is_admin"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new UserModel instance
func NewUser() *UserModel {
    user := &UserModel{
        BaseModel: eloquent.NewBaseModel(),
    }
    
    user.Table("users").
        PrimaryKey("id").
        Fillable("name", "email", "password").
        Hidden("password", "remember_token").
        Casts(map[string]string{
            "email_verified_at": "datetime",
            "is_admin":          "bool",
        })
    
    return user
}

// Define relationships
func (u *UserModel) Posts() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(u)
    return rb.HasMany("posts", "PostModel")
}

func (u *UserModel) Profile() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(u)
    return rb.HasOne("profile", "ProfileModel")
}

// Global static instance for User model - Laravel style
var User = eloquent.NewModelStatic(func() *UserModel {
    return NewUser()
})
```

### 3. Basic Usage

**Laravel-style Model Usage (No Type Assertions Needed!)**

```go
// Find user by email - Direct attribute access!
user, err := models.User.Where("email", "john@example.com").First()
if err != nil {
    log.Fatal(err)
}

// Direct access to typed model attributes - just like Laravel!
fmt.Println("User ID:", user.ID)
fmt.Println("User Name:", user.Name)
fmt.Println("User Email:", user.Email)
fmt.Println("Is Admin:", user.IsAdmin)
fmt.Println("Created At:", user.CreatedAt)

// Get all users with method chaining
users, err := models.User.Where("is_admin", false).
    Where("email_verified_at", "!=", nil).
    OrderBy("created_at", "desc").
    Limit(10).
    Get()

// Loop through typed models - no type assertions!
for _, user := range users {
    fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
}

// Create new user
newUser, err := models.User.Create(map[string]interface{}{
    "name":     "Jane Doe",
    "email":    "jane@example.com",
    "password": "secret123",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created user: %s (ID: %s)\n", newUser.Name, newUser.ID)

// Update user - Method 1: Using Update method with map
err = user.Update(map[string]interface{}{
    "name":  "John Smith Updated",
    "email": "john.smith@example.com",
})
if err != nil {
    log.Fatal(err)
}

// Update user - Method 2: Direct attribute access + Save
user.Name = "John Smith"
user.Email = "john.smith@example.com"
err = user.Save()
if err != nil {
    log.Fatal(err)
}

// Find and update
foundUser, err := models.User.Find(user.ID)
if err != nil {
    log.Fatal(err)
}
foundUser.Name = "Updated Name"
foundUser.Save()

// Delete user
user.Delete() // Soft delete if configured
user.ForceDelete() // Permanent delete
```

## CRUD Operations

Go Eloquent provides comprehensive CRUD (Create, Read, Update, Delete) operations with Laravel-like syntax:

### Create Operations

```go
// Create single record
user, err := models.User.Create(map[string]interface{}{
    "name":     "John Doe",
    "email":    "john@example.com",
    "password": "secret123",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created user: %s (ID: %s)\n", user.Name, user.ID)

// Create with relationships
customer, err := models.Customer.Create(map[string]interface{}{
    "first_name": "Alice",
    "last_name":  "Johnson",
    "email":      "alice@example.com",
    "phone":      "+1234567890",
    "company_id": company.ID, // Link to existing company
})
```

### Read Operations

```go
// Find by primary key
user, err := models.User.Find("123e4567-e89b-12d3-a456-426614174000")

// Find first matching record
user, err := models.User.Where("email", "john@example.com").First()

// Get all records
users, err := models.User.All()

// Get with conditions
activeUsers, err := models.User.Where("status", "active").
    Where("verified", true).
    OrderBy("created_at", "desc").
    Get()

// Check if record exists
exists, err := models.User.Where("email", "john@example.com").Exists()
```

### Update Operations

Go Eloquent provides multiple ways to update records:

```go
// Method 1: Update using Update() method with map
err = user.Update(map[string]interface{}{
    "name":   "John Smith",
    "email":  "john.smith@example.com",
    "status": "premium",
})

// Method 2: Direct attribute modification + Save()
user.Name = "John Updated"
user.Email = "john.updated@example.com"
err = user.Save()

// Method 3: Find and update
foundUser, err := models.User.Find(userID)
if err != nil {
    log.Fatal(err)
}
foundUser.Status = "active"
err = foundUser.Save()

// Method 4: Update multiple records (coming soon)
// affected, err := models.User.Where("status", "inactive").Update(map[string]interface{}{
//     "status": "active",
// })
```

### Delete Operations

```go
// Soft delete (if configured)
err = user.Delete()

// Permanent delete
err = user.ForceDelete()

// Delete by ID
user, err := models.User.Find(userID)
if err != nil {
    log.Fatal(err)
}
err = user.Delete()

// Restore soft deleted record
err = user.Restore()
```

### Complete CRUD Example

```go
package main

import (
    "fmt"
    "log"
    "your-project/models"
    "github.com/crashana/go-eloquent"
)

func main() {
    // Database connection automatically initialized from .env
    defer eloquent.GetManager().CloseAll()

    // CREATE
    user, err := models.User.Create(map[string]interface{}{
        "name":     "John Doe",
        "email":    "john@example.com",
        "password": "secret123",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✅ Created: %s (ID: %s)\n", user.Name, user.ID)

    // READ
    foundUser, err := models.User.Find(user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("📖 Found: %s (%s)\n", foundUser.Name, foundUser.Email)

    // UPDATE
    err = foundUser.Update(map[string]interface{}{
        "name": "John Smith Updated",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("🔄 Updated: %s\n", foundUser.Name)

    // DELETE
    err = foundUser.Delete()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("🗑️  Deleted: %s\n", foundUser.Name)
}
```

**Traditional Query Builder (Still Available)**

```go
// Query users using traditional query builder
db := eloquent.DB()
qb := eloquent.NewQueryBuilder(db)

users, err := qb.Table("users").
    Where("active", true).
    OrderBy("created_at", "desc").
    Limit(10).
    Get()

// Find specific user
user, err := qb.Table("users").Find(1)
```

## Query Builder

The query builder provides a fluent interface for building SQL queries:

```go
db := eloquent.DB()
qb := eloquent.NewQueryBuilder(db)

// Basic queries
users, err := qb.Table("users").
    Select("id", "name", "email").
    Where("active", true).
    Where("created_at", ">=", time.Now().AddDate(0, 0, -30)).
    OrderBy("name", "asc").
    Get()

// Complex queries with joins
posts, err := qb.Table("posts").
    Select("posts.*", "users.name as author_name").
    Join("users", "posts.user_id", "=", "users.id").
    Where("posts.published", true).
    WhereIn("posts.category_id", []interface{}{1, 2, 3}).
    OrderByDesc("posts.created_at").
    Get()

// Aggregates
count, err := qb.Table("users").Count()
sum, err := qb.Table("orders").Sum("amount")
avg, err := qb.Table("products").Avg("price")
```

### Available Query Methods

#### Selecting Data
- `Select(columns...)` - Specify columns to select
- `Distinct()` - Add DISTINCT clause
- `Get()` - Execute and get all results
- `First()` - Get first result
- `Find(id)` - Find by primary key
- `Paginate(page, perPage)` - Paginated results

#### Where Clauses
- `Where(column, operator, value)` - Basic where
- `WhereIn(column, values)` - WHERE IN clause
- `WhereNull(column)` - WHERE NULL clause
- `WhereBetween(column, min, max)` - WHERE BETWEEN clause
- `WhereDate/WhereTime/WhereYear()` - Date-based conditions

#### Joins
- `Join(table, first, operator, second)` - Inner join
- `LeftJoin()` - Left join
- `RightJoin()` - Right join
- `CrossJoin()` - Cross join

#### Ordering & Grouping
- `OrderBy(column, direction)` - Order results
- `GroupBy(columns...)` - Group results
- `Having(column, operator, value)` - Having clause

#### Limiting
- `Limit(count)` / `Take(count)` - Limit results
- `Offset(count)` / `Skip(count)` - Skip results

## Relationships

Define and use relationships just like Eloquent:

```go
// One-to-One
func (u *User) Profile() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(u)
    return rb.HasOne("profile", "Profile")
}

// One-to-Many
func (u *User) Posts() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(u)
    return rb.HasMany("posts", "Post")
}

// Many-to-One
func (p *Post) Author() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(p)
    return rb.BelongsTo("author", "User")
}

// Many-to-Many
func (p *Post) Tags() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(p)
    return rb.BelongsToMany("tags", "Tag", "post_tag")
}

// Polymorphic relationships
func (p *Post) Comments() *eloquent.Relationship {
    rb := eloquent.NewRelationshipBuilder(p)
    return rb.MorphMany("comments", "Comment", "commentable")
}
```

### Relationship Constraints

```go
// Add constraints to relationships
user.Posts().Where("published", true).OrderBy("created_at", "desc").Get()

// Count related models
postsCount, err := user.Posts().Count()

// Check if relationship exists
hasPublishedPosts, err := user.Posts().Where("published", true).Exists()
```

## Scopes

Create reusable query constraints:

```go
// Define scopes
publishedScope := eloquent.PublishedScope()
recentScope := eloquent.RecentScope(30) // Last 30 days
activeScope := eloquent.WhereStatusScope("active")

// Use scopes
qb := eloquent.NewQueryBuilder(db)
eloquent.ApplyScope(qb, publishedScope)
eloquent.ApplyScope(qb, recentScope)

// Chain multiple scopes
combinedScope := eloquent.ChainScopes(publishedScope, recentScope, activeScope)
eloquent.ApplyScope(qb, combinedScope)

// Custom scopes
searchScope := eloquent.SearchScope("john doe", "name", "email")
eloquent.ApplyScope(qb, searchScope)
```

### Available Built-in Scopes

- `PublishedScope()` - Filter published records
- `RecentScope(days)` - Records from last N days
- `SearchScope(query, columns...)` - Text search
- `WhereStatusScope(status)` - Filter by status
- `BetweenDatesScope(start, end)` - Date range filter
- `PaginateScope(page, perPage)` - Pagination
- `OrderScope(column, direction)` - Ordering

## Model Features

### Mass Assignment

```go
user := NewUser()

// Fill multiple attributes
user.Fill(map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
})

// Only fillable attributes are set
user.Fillable("name", "email") // Define fillable fields
user.Guarded("password", "admin") // Define guarded fields
```

### Attribute Casting

```go
user.Casts(map[string]string{
    "email_verified_at": "datetime",
    "is_admin":          "bool",
    "settings":          "json",
    "age":              "int",
})

// Attributes are automatically cast
verified := user.GetAttribute("email_verified_at").(time.Time)
isAdmin := user.GetAttribute("is_admin").(bool)
```

### Hidden/Visible Attributes

```go
// Hide sensitive attributes from serialization
user.Hidden("password", "remember_token")

// Or specify only visible attributes
user.Visible("id", "name", "email")

// Convert to map (respects hidden/visible)
data := user.ToMap()
```

### Timestamps

```go
// Automatic timestamps (default)
user.GetTimestamps() // true

// Disable timestamps
user.WithoutTimestamps()

// Custom timestamp columns
user.GetCreatedAtColumn() // "created_at"
user.GetUpdatedAtColumn() // "updated_at"
```

### Soft Deletes

```go
// Configure soft deletes
user.GetDeletedAtColumn() // "deleted_at"

// Soft delete
user.Delete() // Sets deleted_at timestamp

// Permanent delete
user.ForceDelete()

// Restore soft deleted
user.Restore()

// Query scopes for soft deletes
withTrashed := eloquent.WithTrashedScope()
onlyTrashed := eloquent.OnlyTrashedScope()
```

## Advanced Features

### Transactions

```go
db := eloquent.DB()

err := db.Transaction(func(tx *sqlx.Tx) error {
    // All operations within this function are in a transaction
    
    user := NewUser()
    user.Fill(map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    })
    
    if err := user.Save(); err != nil {
        return err // Transaction will be rolled back
    }
    
    // More operations...
    
    return nil // Transaction will be committed
})
```

### Environment Configuration

```go
// Load custom .env file
err := eloquent.LoadEnv("custom.env")

// Get environment variables with defaults
dbHost := eloquent.Env("DB_HOST", "localhost")
dbPort := eloquent.EnvInt("DB_PORT", 5432)
debugMode := eloquent.EnvBool("DEBUG", false)

// Manual connection with environment variables
err := eloquent.PostgreSQL(eloquent.ConnectionConfig{
    Host:     eloquent.Env("DB_HOST", "localhost"),
    Port:     eloquent.EnvInt("DB_PORT", 5432),
    Database: eloquent.Env("DB_DATABASE"),
    Username: eloquent.Env("DB_USERNAME"),
    Password: eloquent.Env("DB_PASSWORD"),
})
```

### Multiple Connections

```go
// Add named connections
eloquent.GetManager().AddConnection("mysql_main", mysqlConfig)
eloquent.GetManager().AddConnection("postgres_analytics", pgConfig)

// Use specific connection
db := eloquent.DB("mysql_main")
analyticsDB := eloquent.DB("postgres_analytics")
```

### Custom Connection Configuration

```go
config := eloquent.ConnectionConfig{
    Driver:   "mysql",
    Host:     "localhost",
    Port:     3306,
    Database: "myapp",
    Username: "user",
    Password: "password",
    Charset:  "utf8mb4",
    Options: map[string]string{
        "parseTime": "true",
        "loc":       "Local",
    },
}

err := eloquent.GetManager().AddConnection("custom", config)
```

## Examples

Check the [`Examples/`](Examples/) directory for complete working examples:

- [`Examples/main.go`](Examples/main.go) - **Comprehensive CRUD examples** with PostgreSQL
  - ✅ CREATE operations (companies, customers with relationships)
  - ✅ READ operations (queries, Find by ID, method chaining)
  - ✅ UPDATE operations (Update() method, direct attribute access + Save())
  - ✅ DELETE operations (soft delete, permanent delete)
- [`Examples/models/`](Examples/models/) - Model definitions (Company, Customer, Brand)
- [`Examples/.env.example`](Examples/.env.example) - Environment configuration examples
- [`Examples/README.md`](Examples/README.md) - Detailed setup and usage instructions

**Features Demonstrated:**
- Automatic .env database configuration
- Typed models with direct attribute access
- Laravel-style static methods (`Where`, `First`, `All`, `Create`)
- **Complete CRUD operations** (Create, Read, Update, Delete)
- Method chaining with typed returns
- Multiple update methods (Update() with map, direct attribute access + Save())
- Find by ID and update operations
- Automatic UUID generation for PostgreSQL
- Relationship definitions and usage
- No type assertions required!

## API Reference

### Core Types

- `Model` - Base model interface
- `BaseModel` - Default model implementation
- `QueryBuilder` - Fluent query builder
- `Connection` - Database connection wrapper
- `Relationship` - Relationship definition

### Environment & Connection Management

- `LoadEnv(filepath)` - Load environment variables from .env file
- `Env(key, default)` - Get environment variable as string
- `EnvInt(key, default)` - Get environment variable as integer
- `EnvBool(key, default)` - Get environment variable as boolean
- `AutoConnect()` - Automatically connect using .env configuration
- `Init()` - Initialize database connection (alias for AutoConnect)
- `SQLite(database)` - Create SQLite connection
- `MySQL(config)` - Create MySQL connection  
- `PostgreSQL(config)` - Create PostgreSQL connection
- `DB(name...)` - Get database connection
- `GetManager()` - Get connection manager

### Model Static Methods (Laravel-style)

- `models.User.Where(column, value)` - Query with where clause
- `models.User.First()` - Get first record
- `models.User.All()` - Get all records
- `models.User.Get()` - Get records (alias for All)
- `models.User.Find(id)` - Find by primary key
- `models.User.Create(attributes)` - Create new record

### Model Instance Methods

- `Save()` - Save model to database (insert if new, update if exists)
- `Create(attributes)` - Create new record and return typed model
- `Update(attributes)` - Update model attributes using map
- `Delete()` - Delete model (soft delete if configured)
- `ForceDelete()` - Permanently delete model (bypass soft delete)
- `Fill(attributes)` - Mass assign attributes
- `ToMap()` - Convert to map
- `GetAttribute(key)` / `SetAttribute(key, value)` - Attribute access
- `Fresh()` - Reload model from database
- `Refresh()` - Refresh current model instance

### Query Builder Methods

See the Query Builder section above for complete method reference.

## Requirements

- Go 1.21 or higher
- Database drivers:
  - `github.com/go-sql-driver/mysql` (MySQL)
  - `github.com/lib/pq` (PostgreSQL)
  - `github.com/mattn/go-sqlite3` (SQLite)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by Laravel's Eloquent ORM
- Built with [sqlx](https://github.com/jmoiron/sqlx) for database operations
- **Developed by Claude 4 Sonnet AI** - This entire Go Eloquent ORM package was designed and implemented by Anthropic's Claude 4 Sonnet AI

## Roadmap

### ✅ **Completed Features**
- [x] Laravel-style .env configuration
- [x] Typed models with direct attribute access
- [x] Automatic database connection initialization
- [x] Method chaining with typed returns
- [x] PostgreSQL and MySQL support
- [x] **Complete CRUD operations** (Create, Read, Update, Delete)
- [x] Multiple update methods (Update() with map, direct attribute + Save())
- [x] Automatic UUID generation for PostgreSQL
- [x] Find by ID operations
- [x] Basic relationships (HasOne, HasMany, BelongsTo)
- [x] Query builder with fluent interface
- [x] Model static methods (Where, First, All, Create)
- [x] Automatic attribute syncing from database
- [x] Soft delete support
- [x] Mass assignment with fillable/guarded attributes

### 🚧 **In Progress**
- [ ] Advanced relationship features (BelongsToMany, HasManyThrough)
- [ ] Eager loading optimization
- [ ] Query result caching

### 📋 **Planned Features**
- [ ] Schema builder/migrations
- [ ] Model events and observers
- [ ] Database seeding
- [ ] Command-line tools (artisan-like)
- [ ] Performance optimizations
- [ ] SQLite support improvement
- [ ] Advanced query scopes
- [ ] Model factories for testing

---

**Go Eloquent** - Bringing elegant database operations to Go! 🚀 