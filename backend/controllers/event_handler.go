package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

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
