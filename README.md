# Go Eloquent ORM

An Eloquent-inspired ORM for Go that provides the same elegant, expressive syntax for database operations.

## Features

- 🔥 **Eloquent-like API** - Familiar syntax for Eloquent developers
- 🗄️ **Multiple Database Support** - MySQL, PostgreSQL, SQLite
- 🔗 **Relationships** - HasOne, HasMany, BelongsTo, BelongsToMany, and more
- 🔍 **Query Builder** - Fluent, expressive query building
- 🎯 **Scopes** - Reusable query constraints
- 🔄 **Soft Deletes** - Built-in soft delete functionality
- 📝 **Attribute Casting** - Automatic type conversion
- 🚀 **Mass Assignment** - Fillable and guarded attributes
- 📊 **Pagination** - Built-in pagination support
- 🔐 **Transactions** - Database transaction support

## Installation

```bash
go get github.com/crashana/go-eloquent
```

## Quick Start

### 1. Setup Database Connection

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
type User struct {
    *eloquent.BaseModel
}

func NewUser() *User {
    user := &User{
        BaseModel: eloquent.NewBaseModel(),
    }
    
    user.Table("users").
        Fillable("name", "email", "password").
        Hidden("password", "remember_token").
        Casts(map[string]string{
            "email_verified_at": "datetime",
            "is_admin":          "bool",
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
```

### 3. Basic Usage

```go
// Create a new user
user := NewUser()
user.Fill(map[string]interface{}{
    "name":     "John Doe",
    "email":    "john@example.com",
    "password": "secret123",
})
user.Save()

// Query users
db := eloquent.DB()
qb := eloquent.NewQueryBuilder(db)

users, err := qb.Table("users").
    Where("active", true).
    OrderBy("created_at", "desc").
    Limit(10).
    Get()

// Find specific user
user, err := qb.Table("users").Find(1)

// Update user
user.Update(map[string]interface{}{
    "name": "Jane Doe",
})

// Delete user
user.Delete() // Soft delete if configured
user.ForceDelete() // Permanent delete
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

Check the [`example/`](example/) directory for complete working examples:

- [`example/main.go`](example/main.go) - Comprehensive usage examples
- Basic CRUD operations
- Relationship definitions and usage
- Query builder examples
- Scopes and filters
- Model features demonstration

## API Reference

### Core Types

- `Model` - Base model interface
- `BaseModel` - Default model implementation
- `QueryBuilder` - Fluent query builder
- `Connection` - Database connection wrapper
- `Relationship` - Relationship definition

### Connection Management

- `SQLite(database)` - Create SQLite connection
- `MySQL(config)` - Create MySQL connection  
- `PostgreSQL(config)` - Create PostgreSQL connection
- `DB(name...)` - Get database connection
- `GetManager()` - Get connection manager

### Model Methods

- `Save()` - Save model to database
- `Delete()` - Delete model (soft delete if configured)
- `Update(attributes)` - Update model attributes
- `Fill(attributes)` - Mass assign attributes
- `ToMap()` - Convert to map
- `GetAttribute(key)` / `SetAttribute(key, value)` - Attribute access

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

- Inspired by Eloquent ORM
- Built with [sqlx](https://github.com/jmoiron/sqlx) for database operations

## Roadmap

- [ ] Schema builder/migrations
- [ ] Model events and observers
- [ ] Advanced relationship features
- [ ] Query caching
- [ ] Database seeding
- [ ] Command-line tools
- [ ] Performance optimizations

---

**Go Eloquent** - Bringing elegant database operations to Go! 🚀 