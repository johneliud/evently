package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/johneliud/evently/backend/models"
	"github.com/johneliud/evently/backend/repositories"
)

// CalendarHandler handles Google Calendar-related HTTP requests
type CalendarHandler struct {
	CalendarRepo *repositories.CalendarRepository
	EventRepo    *repositories.EventRepository
}

func NewCalendarHandler(calendarRepo *repositories.CalendarRepository, eventRepo *repositories.EventRepository) *CalendarHandler {
	return &CalendarHandler{
		CalendarRepo: calendarRepo,
		EventRepo:    eventRepo,
	}
}

// AuthorizeCalendar initiates the OAuth flow for Google Calendar
func (h *CalendarHandler) AuthorizeCalendar(w http.ResponseWriter, r *http.Request) {
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

	// Generate a state token to prevent request forgery
	state := fmt.Sprintf("user-%d-%d", userID, time.Now().Unix())

	// Get the authorization URL
	authURL := h.CalendarRepo.GetAuthURL(state)

	// Return the authorization URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"auth_url": authURL,
	})
}

// CalendarCallback handles the OAuth callback from Google
func (h *CalendarHandler) CalendarCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get the authorization code and state from the query parameters
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		log.Println("Missing authorization code")
		return
	}

	// Extract user ID from state
	var userID int
	_, err := fmt.Sscanf(state, "user-%d-", &userID)
	if err != nil || userID == 0 {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		log.Printf("Invalid state parameter: %v\n", err)
		return
	}

	// Exchange the authorization code for a token
	ctx := context.Background()
	token, err := h.CalendarRepo.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		log.Printf("Failed to exchange authorization code: %v\n", err)
		return
	}

	// Store the token for the user
	err = h.CalendarRepo.StoreUserToken(userID, token)
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		log.Printf("Failed to store token: %v\n", err)
		return
	}

	var url string

	environmentState := os.Getenv("ENVIRONMENT")
	if environmentState == "production" {
		url = "https://evently-dgq9.onrender.com"
	} else {
		url = "http://localhost:5173"
	}

	// Redirect to the frontend success page
	http.Redirect(w, r, url+"/calendar-connected", http.StatusFound)
}

// AddEventToCalendar adds an event to the user's Google Calendar
func (h *CalendarHandler) AddEventToCalendar(w http.ResponseWriter, r *http.Request) {
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

	// Get event ID from request body
	var req struct {
		EventID int `json:"event_id"`
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var stringReq struct {
			EventID string `json:"event_id"`
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if jsonErr := json.NewDecoder(r.Body).Decode(&stringReq); jsonErr == nil {

			eventID, convErr := strconv.Atoi(stringReq.EventID)
			if convErr == nil {
				req.EventID = eventID
			} else {
				http.Error(w, "Invalid event ID format", http.StatusBadRequest)
				log.Printf("Invalid event ID format: %v\n", convErr)
				return
			}
		} else {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Printf("Invalid request body: %v\n", err)
			return
		}
	}

	// Get the event from the database
	event, err := h.EventRepo.GetEventByID(req.EventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		log.Printf("Event not found: %v\n", err)
		return
	}

	// Get the user's token
	token, err := h.CalendarRepo.GetUserToken(userID)
	if err != nil {
		// If the token is not found, redirect to authorization
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Google Calendar authorization required",
			"status":  "authorization_required",
		})
		return
	}

	// Check if the token is expired and refresh if needed
	if token.Expiry.Before(time.Now()) {
		ctx := context.Background()
		newToken, err := h.CalendarRepo.RefreshToken(ctx, token, userID)
		if err != nil {
			http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
			log.Printf("Failed to refresh token: %v\n", err)
			return
		}
		token = newToken
	}

	// Convert EventWithOrganizer to Event
	eventModel := &models.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: event.Description,
		Date:        event.Date,
		Location:    event.Location,
		UserID:      event.UserID,
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}

	// Add the event to Google Calendar
	ctx := context.Background()
	calendarEvent, err := h.CalendarRepo.AddEventToCalendar(ctx, token, eventModel)
	if err != nil {
		http.Error(w, "Failed to add event to Google Calendar", http.StatusInternalServerError)
		log.Printf("Failed to add event to Google Calendar: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Event added to Google Calendar successfully",
		"calendar_event": calendarEvent,
	})
}

// CheckCalendarConnection checks if the user has connected their Google Calendar
func (h *CalendarHandler) CheckCalendarConnection(w http.ResponseWriter, r *http.Request) {
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

	// Check if the user has a token
	token, err := h.CalendarRepo.GetUserToken(userID)
	if err != nil {
		// Token not found, user needs to connect
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"connected": false,
		})
		return
	}

	// Check if the token is valid
	ctx := context.Background()
	_, err = h.CalendarRepo.GetCalendarService(ctx, token)
	if err != nil {
		// Token is invalid, user needs to reconnect
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"connected": false,
		})
		return
	}

	// Token is valid, user is connected
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"connected": true,
	})
}
