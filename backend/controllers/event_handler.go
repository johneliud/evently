package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/johneliud/evently/backend/models"
	"github.com/johneliud/evently/backend/repositories"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	EventRepo *repositories.EventRepository
}

func NewEventHandler(eventRepo *repositories.EventRepository) *EventHandler {
	return &EventHandler{EventRepo: eventRepo}
}

// CreateEvent handles event creation
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized: %v\n", err)
		return
	}

	var req models.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v\n", err)
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Location) == "" {
		http.Error(w, "Title and location are required", http.StatusBadRequest)
		log.Println("Title and location are required")
		return
	}

	// Create event
	id, err := h.EventRepo.CreateEvent(req, userID)
	if err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		log.Printf("Failed to create event: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Event created successfully",
	})
	log.Println("Event created successfully")
}

// GetUserEvents handles retrieving events for a user
func (h *EventHandler) GetUserEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized: %v\n", err)
		return
	}

	// Get events
	events, err := h.EventRepo.GetEventsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get events", http.StatusInternalServerError)
		log.Printf("Failed to get events: %v\n", err)
		return
	}

	if events == nil {
		events = []models.Event{}
	}

	// Return events
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// GetUpcomingEvents handles retrieving upcoming public events
func (h *EventHandler) GetUpcomingEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get upcoming events
	events, err := h.EventRepo.GetUpcomingEvents()
	if err != nil {
		http.Error(w, "Failed to get upcoming events", http.StatusInternalServerError)
		log.Printf("Failed to get upcoming events: %v\n", err)
		return
	}

	// Ensure we return an empty array instead of null if no events are found
	if events == nil {
		events = []models.Event{}
	}

	// Return events
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// GetEventByID handles retrieving a single event by ID
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Extract event ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	idStr := segments[len(segments)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Get event by ID
	event, err := h.EventRepo.GetEventByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
			log.Printf("Event not found: %v\n", err)
			return
		}
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		log.Printf("Failed to get event: %v\n", err)
		return
	}

	// Return event
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DeleteEvent handles event deletion
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized: %v\n", err)
		return
	}

	// Extract event ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	idStr := segments[len(segments)-1]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Check if the event belongs to the user
	event, err := h.EventRepo.GetEventByID(eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
			log.Printf("Event not found: %v\n", err)
			return
		}
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		log.Printf("Failed to get event: %v\n", err)
		return
	}

	if event.UserID != userID {
		http.Error(w, "Unauthorized: You can only delete your own events", http.StatusForbidden)
		log.Printf("Unauthorized: User %d attempted to delete event %d owned by user %d\n", userID, eventID, event.UserID)
		return
	}

	// Delete the event
	err = h.EventRepo.DeleteEvent(eventID)
	if err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		log.Printf("Failed to delete event: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Event deleted successfully",
	})
	log.Printf("Event %d deleted successfully by user %d\n", eventID, userID)
}

// UpdateEvent handles event updates
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get user ID from token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized: %v\n", err)
		return
	}

	// Extract event ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	idStr := segments[len(segments)-1]
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Check if the event belongs to the user
	event, err := h.EventRepo.GetEventByID(eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
			log.Printf("Event not found: %v\n", err)
			return
		}
		http.Error(w, "Failed to get event", http.StatusInternalServerError)
		log.Printf("Failed to get event: %v\n", err)
		return
	}

	if event.UserID != userID {
		http.Error(w, "Unauthorized: You can only update your own events", http.StatusForbidden)
		log.Printf("Unauthorized: User %d attempted to update event %d owned by user %d\n", userID, eventID, event.UserID)
		return
	}

	// Parse request body
	var req models.EventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v\n", err)
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Location) == "" {
		http.Error(w, "Title and location are required", http.StatusBadRequest)
		log.Println("Title and location are required")
		return
	}

	// Update the event
	err = h.EventRepo.UpdateEvent(eventID, req)
	if err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		log.Printf("Failed to update event: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Event updated successfully",
	})
	log.Printf("Event %d updated successfully by user %d\n", eventID, userID)
}

// SearchEvents handles searching and filtering events
func (h *EventHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Parse query parameters
	query := r.URL.Query().Get("q")
	location := r.URL.Query().Get("location")
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	var startDate, endDate *time.Time

	// Parse start date if provided
	if startDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start date format. Use YYYY-MM-DD", http.StatusBadRequest)
			log.Printf("Invalid start date format: %v\n", err)
			return
		}
		startDate = &parsedDate
	}

	// Parse end date if provided
	if endDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end date format. Use YYYY-MM-DD", http.StatusBadRequest)
			log.Printf("Invalid end date format: %v\n", err)
			return
		}
		// Set the end date to the end of the day
		parsedDate = parsedDate.Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)
		endDate = &parsedDate
	}

	// Search events
	events, err := h.EventRepo.SearchEvents(query, location, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to search events", http.StatusInternalServerError)
		log.Printf("Failed to search events: %v\n", err)
		return
	}

	// Ensure we return an empty array instead of null if no events are found
	if events == nil {
		events = []models.EventWithOrganizer{}
	}

	// Return events
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Helper function to extract user ID from JWT token
func getUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err := errors.New("authorization header is required")
		log.Printf("Authorization header is required: %v\n", err)
		return 0, err
	}

	// Remove 'Bearer ' prefix
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %v\n", err)
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert user_id to int
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return 0, jwt.ErrTokenInvalidClaims
		}
		userID := int(userIDFloat)
		return userID, nil
	}

	return 0, jwt.ErrTokenInvalidClaims
}
