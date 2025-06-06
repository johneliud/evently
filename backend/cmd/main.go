package main

import (
	"log"
	"os"

	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if it exists
	// This will be skipped in production where environment variables are set differently
	_ = godotenv.Load(".env")

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

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Start the server
	if err := srv.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
