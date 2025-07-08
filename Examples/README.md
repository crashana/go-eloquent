# Go Eloquent Examples

This directory contains examples showing how to use Go Eloquent ORM with automatic database configuration.

## Quick Start

### 1. Configure Database

Copy the example environment file and customize it:

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```env
# PostgreSQL Configuration
DB_CONNECTION=pgsql
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=testApp
DB_USERNAME=postgres
DB_PASSWORD=postgres

# Or MySQL Configuration
# DB_CONNECTION=mysql
# DB_HOST=localhost
# DB_PORT=3306
# DB_DATABASE=testApp
# DB_USERNAME=root
# DB_PASSWORD=password
```

### 2. Run the Example

```bash
go run main.go
```

That's it! The database connection is automatically configured from your `.env` file.

## Supported Database Types

| DB_CONNECTION | Description |
|---------------|-------------|
| `pgsql` | PostgreSQL |
| `postgres` | PostgreSQL (alias) |
| `postgresql` | PostgreSQL (alias) |
| `mysql` | MySQL |

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DB_CONNECTION` | Database type | No | `pgsql` |
| `DB_HOST` | Database host | No | `localhost` |
| `DB_PORT` | Database port | No | `5432` (pgsql) or `3306` (mysql) |
| `DB_DATABASE` | Database name | **Yes** | - |
| `DB_USERNAME` | Database username | **Yes** | - |
| `DB_PASSWORD` | Database password | No | - |
| `DB_CHARSET` | Database charset (MySQL only) | No | `utf8mb4` |

## Features Demonstrated

The example demonstrates:

1. **Automatic Database Connection** - No manual configuration needed
2. **Typed Models** - Direct access to model attributes without type assertions
3. **Laravel-like Syntax** - Familiar query builder methods
4. **Method Chaining** - Fluent query building with typed returns
5. **Relationships** - Model relationships (HasMany, BelongsTo, etc.)

## Laravel vs Go Eloquent

### Laravel (PHP)
```php
// .env file
DB_CONNECTION=pgsql
DB_HOST=localhost
DB_DATABASE=testApp
DB_USERNAME=postgres
DB_PASSWORD=postgres

// PHP code
$company = Company::where('name', 'like', '%test%')->first();
echo $company->name; // Direct access
```

### Go Eloquent
```env
# .env file
DB_CONNECTION=pgsql
DB_HOST=localhost
DB_DATABASE=testApp
DB_USERNAME=postgres
DB_PASSWORD=postgres
```

```go
// Go code
company, err := models.Company.Where("name", "like", "%test%").First()
fmt.Println(company.Name) // Direct access - no type assertions!
```

## Manual Configuration (Alternative)

If you prefer not to use `.env` files, you can still configure manually:

```go
err := eloquent.PostgreSQL(eloquent.ConnectionConfig{
    Host:     "localhost",
    Port:     5432,
    Database: "testApp",
    Username: "postgres",
    Password: "postgres",
})
```

## Models

The example uses models from the `models/` package:

- `CompanyModel` - Demonstrates basic model usage
- `CustomerModel` - Shows relationships and attributes
- `BrandModel` - Additional model example

Each model provides:
- Direct attribute access (no type assertions)
- Automatic attribute syncing from database
- Laravel-like static methods (`Where`, `First`, `All`, etc.)
- Method chaining with typed returns 