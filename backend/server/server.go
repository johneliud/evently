package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/johneliud/evently/backend/controllers"
	"github.com/johneliud/evently/backend/repositories"
	"github.com/johneliud/evently/backend/services"
)

// Server represents the HTTP server
type Server struct {
	Database     *sql.DB
	Services     *ServiceContainer
	Repositories *RepositoryContainer
	Handlers     *HandlerContainer
	Mux          *http.ServeMux
}

// ServiceContainer holds all services
type ServiceContainer struct {
	EmailService *services.EmailService
}

// RepositoryContainer holds all repositories
type RepositoryContainer struct {
	UserRepo     *repositories.UserRepository
	EventRepo    *repositories.EventRepository
	RSVPRepo     *repositories.RSVPRepository
	CalendarRepo *repositories.CalendarRepository
}

// HandlerContainer holds all handlers
type HandlerContainer struct {
	UserHandler     *controllers.UserHandler
	EventHandler    *controllers.EventHandler
	RSVPHandler     *controllers.RSVPHandler
	CalendarHandler *controllers.CalendarHandler
}

// NewServer creates a new server instance
func NewServer(database *sql.DB) (*Server, error) {
	server := &Server{
		Database: database,
		Mux:      http.NewServeMux(),
	}

	// Initialize services and repositories
	if err := server.initServicesAndRepositories(); err != nil {
		return nil, err
	}

	// Initialize handlers
	server.initHandlers()

	// Setup routes
	server.setupRoutes()

	return server, nil
}

// initServicesAndRepositories initializes all services and repositories
func (s *Server) initServicesAndRepositories() error {
	// Initialize services
	emailService := services.NewEmailService()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(s.Database)
	eventRepo := repositories.NewEventRepository(s.Database)
	rsvpRepo := repositories.NewRSVPRepository(s.Database)

	// Initialize Google Calendar repository
	calendarRepo, err := repositories.NewCalendarRepository()
	if err != nil {
		return fmt.Errorf("failed to initialize calendar repository: %v", err)
	}

	s.Services = &ServiceContainer{
		EmailService: emailService,
	}

	s.Repositories = &RepositoryContainer{
		UserRepo:     userRepo,
		EventRepo:    eventRepo,
		RSVPRepo:     rsvpRepo,
		CalendarRepo: calendarRepo,
	}

	return nil
}

// initHandlers initializes all handlers
func (s *Server) initHandlers() {
	s.Handlers = &HandlerContainer{
		UserHandler:     controllers.NewUserHandler(s.Repositories.UserRepo),
		EventHandler:    controllers.NewEventHandler(s.Repositories.EventRepo),
		RSVPHandler:     controllers.NewRSVPHandler(s.Repositories.RSVPRepo, s.Repositories.EventRepo, s.Repositories.UserRepo, s.Services.EmailService),
		CalendarHandler: controllers.NewCalendarHandler(s.Repositories.CalendarRepo, s.Repositories.EventRepo),
	}
}

// setupRoutes sets up all API routes
func (s *Server) setupRoutes() {
	// Register handlers
	s.Mux.Handle("/api/signup", corsMiddleware(http.HandlerFunc(s.Handlers.UserHandler.SignUp)))
	s.Mux.Handle("/api/signin", corsMiddleware(http.HandlerFunc(s.Handlers.UserHandler.SignIn)))

	// Add Google OAuth routes
	s.Mux.Handle("/api/auth/google", corsMiddleware(http.HandlerFunc(s.Handlers.UserHandler.GoogleAuthURL)))
	s.Mux.Handle("/api/auth/google/callback", corsMiddleware(http.HandlerFunc(s.Handlers.UserHandler.GoogleCallback)))

	// Event routes
	s.Mux.Handle("/api/events", corsMiddleware(http.HandlerFunc(s.Handlers.EventHandler.CreateEvent)))
	s.Mux.Handle("/api/events/user", corsMiddleware(http.HandlerFunc(s.Handlers.EventHandler.GetUserEvents)))
	s.Mux.Handle("/api/events/upcoming", corsMiddleware(http.HandlerFunc(s.Handlers.EventHandler.GetUpcomingEvents)))
	s.Mux.Handle("/api/events/search", corsMiddleware(http.HandlerFunc(s.Handlers.EventHandler.SearchEvents)))

	// Google Calendar endpoints
	s.Mux.Handle("/api/calendar/authorize", corsMiddleware(http.HandlerFunc(s.Handlers.CalendarHandler.AuthorizeCalendar)))
	s.Mux.Handle("/api/calendar/callback", corsMiddleware(http.HandlerFunc(s.Handlers.CalendarHandler.CalendarCallback)))
	s.Mux.Handle("/api/calendar/add-event", corsMiddleware(http.HandlerFunc(s.Handlers.CalendarHandler.AddEventToCalendar)))
	s.Mux.Handle("/api/calendar/check-connection", corsMiddleware(http.HandlerFunc(s.Handlers.CalendarHandler.CheckCalendarConnection)))

	// Dynamic event routes
	s.Mux.Handle("/api/events/", corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/rsvp") {
			switch r.Method {
			case http.MethodGet:
				s.Handlers.RSVPHandler.GetRSVP(w, r)
			case http.MethodPost, http.MethodPut:
				s.Handlers.RSVPHandler.CreateOrUpdateRSVP(w, r)
			case http.MethodDelete:
				s.Handlers.RSVPHandler.DeleteRSVP(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else if strings.HasSuffix(path, "/rsvp/count") {
			s.Handlers.RSVPHandler.GetRSVPCount(w, r)
		} else if strings.HasSuffix(path, "/rsvps") {
			s.Handlers.RSVPHandler.GetRSVPs(w, r)
		} else {
			switch r.Method {
			case http.MethodGet:
				s.Handlers.EventHandler.GetEventByID(w, r)
			case http.MethodDelete:
				s.Handlers.EventHandler.DeleteEvent(w, r)
			case http.MethodPut:
				s.Handlers.EventHandler.UpdateEvent(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})))
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	fmt.Printf("Server starting on %s\n", addr)
	return http.ListenAndServe(addr, corsMiddleware(s.Mux))
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")

		// If no origin is provided, use the default frontend URL based on environment
		if origin == "" {
			// Check if we're in production
			if os.Getenv("ENVIRONMENT") == "production" {
				origin = os.Getenv("FRONTEND_URL")
				if origin == "" {
					origin = "https://evently-dgq9.onrender.com"
				}
			} else {
				origin = "http://localhost:5173"
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
