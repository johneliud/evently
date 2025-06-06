package main

import (
	"log"
	"os"

	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/server"
	"github.com/joho/godotenv"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Load environment variables from .env file if it exists
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Failed loading .env: %v\n", err)
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

	log.Println("Successfully connected to database and ran migrations")

	// Create and initialize the server
	srv, err := server.NewServer(database)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Start the server
	environmentState := os.Getenv("ENVIRONMENT")
	if environmentState == "production" {
		if err := srv.Start("0.0.0.0:" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
		return
	}

	if err := srv.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
