package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/johneliud/evently/backend/controllers"
	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/repositories"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
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
	eventRepo := repositories.NewEventRepository(database)
	rsvpRepo := repositories.NewRSVPRepository(database)

	// Initialize handlers
	userHandler := controllers.NewUserHandler(userRepo)
	eventHandler := controllers.NewEventHandler(eventRepo)
	rsvpHandler := controllers.NewRSVPHandler(rsvpRepo, eventRepo)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.Handle("/api/signup", corsMiddleware(http.HandlerFunc(userHandler.SignUp)))
	mux.Handle("/api/signin", corsMiddleware(http.HandlerFunc(userHandler.SignIn)))
	mux.Handle("/api/events", corsMiddleware(http.HandlerFunc(eventHandler.CreateEvent)))
	mux.Handle("/api/events/user", corsMiddleware(http.HandlerFunc(eventHandler.GetUserEvents)))
	mux.Handle("/api/events/upcoming", corsMiddleware(http.HandlerFunc(eventHandler.GetUpcomingEvents)))
	mux.Handle("/api/events/search", corsMiddleware(http.HandlerFunc(eventHandler.SearchEvents)))
	mux.Handle("/api/events/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/rsvp") {
			switch r.Method {
			case http.MethodGet:
				rsvpHandler.GetRSVP(w, r)
			case http.MethodPost, http.MethodPut:
				rsvpHandler.CreateOrUpdateRSVP(w, r)
			case http.MethodDelete:
				rsvpHandler.DeleteRSVP(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if strings.HasSuffix(path, "/rsvp/count") {
			rsvpHandler.GetRSVPCount(w, r)
		} else if strings.HasSuffix(path, "/rsvps") {
			rsvpHandler.GetRSVPs(w, r)
		} else {
			switch r.Method {
			case http.MethodGet:
				eventHandler.GetEventByID(w, r)
			case http.MethodDelete:
				eventHandler.DeleteEvent(w, r)
			case http.MethodPut:
				eventHandler.UpdateEvent(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})))

	// Start server
	fmt.Println("Server starting on :9000")
	if err := http.ListenAndServe(":9000", corsMiddleware(mux)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware adds CORS headers to all responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

		// Send 200 OK response for preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
