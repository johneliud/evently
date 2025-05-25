package main

import (
	"log"

	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Successfully connected to database")

	// Create and initialize the server
	srv, err := server.NewServer(database)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start the server
	if err := srv.Start(":9000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
