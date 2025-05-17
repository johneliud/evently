package db

import (
	"database/sql"
)

// RunMigrations creates necessary tables if they don't exist
func RunMigrations(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            confirmed_password VARCHAR(255) NOT NULL,
            first_name VARCHAR(100) NOT NULL,
            last_name VARCHAR(100) NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )
    `)
	return err
}
