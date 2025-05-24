package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/johneliud/evently/backend/models"
	"github.com/johneliud/evently/backend/repositories"
)

// RSVPHandler handles RSVP-related HTTP requests
type RSVPHandler struct {
	RSVPRepo  *repositories.RSVPRepository
	EventRepo *repositories.EventRepository
}

func NewRSVPHandler(rsvpRepo *repositories.RSVPRepository, eventRepo *repositories.EventRepository) *RSVPHandler {
	return &RSVPHandler{
		RSVPRepo:  rsvpRepo,
		EventRepo: eventRepo,
	}
}

// CreateOrUpdateRSVP handles creating or updating an RSVP
func (h *RSVPHandler) CreateOrUpdateRSVP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
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
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	eventIDStr := segments[len(segments)-2]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Check if event exists
	_, err = h.EventRepo.GetEventByID(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		log.Printf("Event not found: %v\n", err)
		return
	}

	// Parse request body
	var req models.RSVPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v\n", err)
		return
	}

	// Validate status
	if req.Status != "going" && req.Status != "maybe" && req.Status != "not_going" {
		http.Error(w, "Invalid status. Must be 'going', 'maybe', or 'not_going'", http.StatusBadRequest)
		log.Printf("Invalid status: %s\n", req.Status)
		return
	}

	// Create or update RSVP
	err = h.RSVPRepo.CreateOrUpdateRSVP(eventID, userID, req.Status)
	if err != nil {
		http.Error(w, "Failed to create/update RSVP", http.StatusInternalServerError)
		log.Printf("Failed to create/update RSVP: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "RSVP updated successfully",
	})
	log.Printf("RSVP updated successfully for event %d by user %d with status %s\n", eventID, userID, req.Status)
}

// GetRSVP handles retrieving a user's RSVP for an event
func (h *RSVPHandler) GetRSVP(w http.ResponseWriter, r *http.Request) {
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

	// Extract event ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	eventIDStr := segments[len(segments)-2]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Get RSVP
	rsvp, err := h.RSVPRepo.GetRSVPByEventAndUser(eventID, userID)
	if err != nil {
		http.Error(w, "Failed to get RSVP", http.StatusInternalServerError)
		log.Printf("Failed to get RSVP: %v\n", err)
		return
	}

	// Return RSVP (or null if not found)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rsvp)
}

// DeleteRSVP handles deleting a user's RSVP for an event
func (h *RSVPHandler) DeleteRSVP(w http.ResponseWriter, r *http.Request) {
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
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	eventIDStr := segments[len(segments)-2]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Delete RSVP
	err = h.RSVPRepo.DeleteRSVP(eventID, userID)
	if err != nil {
		http.Error(w, "Failed to delete RSVP", http.StatusInternalServerError)
		log.Printf("Failed to delete RSVP: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "RSVP deleted successfully",
	})
	log.Printf("RSVP deleted successfully for event %d by user %d\n", eventID, userID)
}

// GetRSVPCount handles retrieving the count of RSVPs by status for an event
func (h *RSVPHandler) GetRSVPCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Extract event ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	eventIDStr := segments[len(segments)-3]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Get RSVP count
	count, err := h.RSVPRepo.GetRSVPCountByEvent(eventID)
	if err != nil {
		http.Error(w, "Failed to get RSVP count", http.StatusInternalServerError)
		log.Printf("Failed to get RSVP count: %v\n", err)
		return
	}

	// Return count
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
}

// GetRSVPs handles retrieving all RSVPs for an event
func (h *RSVPHandler) GetRSVPs(w http.ResponseWriter, r *http.Request) {
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

	// Extract event ID from URL path
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		log.Println("Invalid URL")
		return
	}

	eventIDStr := segments[len(segments)-2]
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		log.Printf("Invalid event ID: %v\n", err)
		return
	}

	// Check if the user is the event creator
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

	// Only the event creator can see the full list of RSVPs
	if event.UserID != userID {
		http.Error(w, "Unauthorized. Only the event creator can view the attendee list", http.StatusForbidden)
		log.Printf("Unauthorized access to RSVPs: User %d tried to access RSVPs for event %d created by user %d\n", userID, eventID, event.UserID)
		return
	}

	// Get RSVPs
	rsvps, err := h.RSVPRepo.GetRSVPsByEvent(eventID)
	if err != nil {
		http.Error(w, "Failed to get RSVPs", http.StatusInternalServerError)
		log.Printf("Failed to get RSVPs: %v\n", err)
		return
	}

	// Return RSVPs
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rsvps)
	log.Printf("RSVPs retrieved successfully for event %d by creator %d\n", eventID, userID)
}
