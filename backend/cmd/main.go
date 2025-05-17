package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johneliud/evently/backend/controllers"
	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/repositories"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	fmt.Println("Successfully connected to database")

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)

	// Initialize controllers
	userHandler := controllers.NewUserHandler(userRepo)

	// Set up routes
	http.HandleFunc("/api/signup", userHandler.SignUp)

	// Enable CORS
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		http.DefaultServeMux.ServeHTTP(w, r)
	})

	// Start server
	fmt.Println("Server starting on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
