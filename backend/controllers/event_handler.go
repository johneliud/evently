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

