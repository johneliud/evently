package db

import (
	"database/sql"
	"fmt"

	"github.com/johneliud/evently/backend/config"
	_ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database
func Connect() (*sql.DB, error) {
	config := config.GetDatabaseConfig()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
