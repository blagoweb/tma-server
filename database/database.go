package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func New(databaseURL string) (*Database, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Initialize database tables
	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func initTables(db *sql.DB) error {
	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		telegram_id BIGINT UNIQUE NOT NULL,
		username VARCHAR(255),
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create user_sessions table for JWT tokens
	createSessionsTable := `
	CREATE TABLE IF NOT EXISTS user_sessions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		token_hash VARCHAR(255) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create pages table (example for your app)
	createPagesTable := `
	CREATE TABLE IF NOT EXISTS pages (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		json_data JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	tables := []string{createUsersTable, createSessionsTable, createPagesTable}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("Database tables initialized successfully")
	return nil
}

func runMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Migration 1: Add json_data column to pages table if it doesn't exist
	addJSONDataColumn := `
	DO $$ 
	BEGIN 
		-- Check if pages table exists
		IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'pages') THEN
			-- Check if json_data column exists
			IF NOT EXISTS (
				SELECT 1 FROM information_schema.columns 
				WHERE table_name = 'pages' AND column_name = 'json_data'
			) THEN 
				ALTER TABLE pages ADD COLUMN json_data JSONB;
				RAISE NOTICE 'Added json_data column to pages table';
			ELSE 
				RAISE NOTICE 'json_data column already exists in pages table';
			END IF;
		ELSE
			RAISE NOTICE 'pages table does not exist, will be created by initTables';
		END IF;
	END $$;
	`

	if _, err := db.Exec(addJSONDataColumn); err != nil {
		return fmt.Errorf("failed to add json_data column: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
