package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/johneliud/evently/backend/controllers"
	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/repositories"
	"github.com/johneliud/evently/backend/services"
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

	// Initialize email service
	emailService := services.NewEmailService()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)
	eventRepo := repositories.NewEventRepository(database)
	rsvpRepo := repositories.NewRSVPRepository(database)

	// Initialize Google Calendar repository
	calendarRepo, err := repositories.NewCalendarRepository()
	if err != nil {
		log.Fatalf("Failed to initialize calendar repository: %v", err)
	}

	// Initialize handlers
	userHandler := controllers.NewUserHandler(userRepo)
	eventHandler := controllers.NewEventHandler(eventRepo)
	rsvpHandler := controllers.NewRSVPHandler(rsvpRepo, eventRepo, userRepo, emailService)
	calendarHandler := controllers.NewCalendarHandler(calendarRepo, eventRepo)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.Handle("/api/signup", corsMiddleware(http.HandlerFunc(userHandler.SignUp)))
	mux.Handle("/api/signin", corsMiddleware(http.HandlerFunc(userHandler.SignIn)))
	// Add Google OAuth routes
	mux.Handle("/api/auth/google", corsMiddleware(http.HandlerFunc(userHandler.GoogleAuthURL)))
	mux.Handle("/api/auth/google/callback", corsMiddleware(http.HandlerFunc(userHandler.GoogleCallback)))
	mux.Handle("/api/events", corsMiddleware(http.HandlerFunc(eventHandler.CreateEvent)))
	mux.Handle("/api/events/user", corsMiddleware(http.HandlerFunc(eventHandler.GetUserEvents)))
	mux.Handle("/api/events/upcoming", corsMiddleware(http.HandlerFunc(eventHandler.GetUpcomingEvents)))
	mux.Handle("/api/events/search", corsMiddleware(http.HandlerFunc(eventHandler.SearchEvents)))

	// Google Calendar endpoints
	mux.Handle("/api/calendar/authorize", corsMiddleware(http.HandlerFunc(calendarHandler.AuthorizeCalendar)))
	mux.Handle("/api/calendar/callback", corsMiddleware(http.HandlerFunc(calendarHandler.CalendarCallback)))
	mux.Handle("/api/calendar/add-event", corsMiddleware(http.HandlerFunc(calendarHandler.AddEventToCalendar)))
	mux.Handle("/api/calendar/check-connection", corsMiddleware(http.HandlerFunc(calendarHandler.CheckCalendarConnection)))
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

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:5173" // Default to frontend URL if Origin header is not set
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Allow cookies

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
