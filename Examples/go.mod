module postgres_test

go 1.21

replace github.com/crashana/go-eloquent => ../

require github.com/crashana/go-eloquent v0.0.0-00010101000000-000000000000

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
)
