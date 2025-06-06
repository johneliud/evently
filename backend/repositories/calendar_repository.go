package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/johneliud/evently/backend/db"
	"github.com/johneliud/evently/backend/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// CalendarRepository handles Google Calendar operations
type CalendarRepository struct {
	Config *oauth2.Config
}

func NewCalendarRepository() (*CalendarRepository, error) {
	credFile := os.Getenv("GOOGLE_CREDENTIALS_FILE")
	if credFile == "" {
		credFile = "google_client_credentials.json"
	}

	data, err := os.ReadFile(credFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client credentials file: %v", err)
	}
	credentialsJSON := string(data)

	// Parse the credentials
	config, err := google.ConfigFromJSON([]byte(credentialsJSON), calendar.CalendarEventsScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client credentials: %v", err)
	}

	var url string

	environmentState := os.Getenv("ENVIRONMENT")
	if environmentState == "production" {
		url = "https://evently-backend-gs5n.onrender.com"
	} else {
		url = "http://localhost:9000"
	}

	redirectURI := url + "/api/calendar/callback"
	config.RedirectURL = redirectURI

	return &CalendarRepository{
		Config: config,
	}, nil
}

// GetAuthURL returns the URL to visit to authorize the application
func (r *CalendarRepository) GetAuthURL(state string) string {
	return r.Config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// Exchange exchanges an authorization code for a token
func (r *CalendarRepository) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return r.Config.Exchange(ctx, code)
}

// GetCalendarService returns a Google Calendar service
func (r *CalendarRepository) GetCalendarService(ctx context.Context, token *oauth2.Token) (*calendar.Service, error) {
	client := r.Config.Client(ctx, token)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create calendar service: %v", err)
	}
	return srv, nil
}

// AddEventToCalendar adds an event to the user's Google Calendar
func (r *CalendarRepository) AddEventToCalendar(ctx context.Context, token *oauth2.Token, event *models.Event) (*calendar.Event, error) {
	// Get calendar service
	srv, err := r.GetCalendarService(ctx, token)
	if err != nil {
		return nil, err
	}

	// Create calendar event
	calendarEvent := &calendar.Event{
		Summary:     event.Title,
		Description: event.Description,
		Start: &calendar.EventDateTime{
			DateTime: event.Date.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			// Add 2 hours by default if no end time is specified
			DateTime: event.Date.Add(2 * time.Hour).Format(time.RFC3339),
			TimeZone: "UTC",
		},
		Location: event.Location,
	}

	// Insert the event
	calendarEvent, err = srv.Events.Insert("primary", calendarEvent).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create event in calendar: %v", err)
	}

	return calendarEvent, nil
}

// StoreUserToken stores a user's OAuth token in the database
func (r *CalendarRepository) StoreUserToken(userID int, token *oauth2.Token) error {
	// Convert token to JSON
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("unable to marshal token: %v", err)
	}

	// Get a database connection
	db, err := db.Connect()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	// Check if a token already exists for this user
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_calendar_tokens WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if token exists: %v", err)
	}

	// Either insert a new token or update the existing one
	if exists {
		_, err = db.Exec(
			"UPDATE user_calendar_tokens SET token = $1, updated_at = NOW() WHERE user_id = $2",
			string(tokenJSON), userID,
		)
	} else {
		_, err = db.Exec(
			"INSERT INTO user_calendar_tokens (user_id, token, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())",
			userID, string(tokenJSON),
		)
	}

	if err != nil {
		return fmt.Errorf("unable to store token in database: %v", err)
	}

	return nil
}

// GetUserToken retrieves a user's OAuth token from the database
func (r *CalendarRepository) GetUserToken(userID int) (*oauth2.Token, error) {
	// Get a database connection
	db, err := db.Connect()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	// Retrieve the token from the database
	var tokenJSON string
	err = db.QueryRow("SELECT token FROM user_calendar_tokens WHERE user_id = $1", userID).Scan(&tokenJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no token found for user: %v", userID)
		}
		return nil, fmt.Errorf("unable to retrieve token from database: %v", err)
	}

	// Parse the token
	var token oauth2.Token
	if err := json.Unmarshal([]byte(tokenJSON), &token); err != nil {
		return nil, fmt.Errorf("unable to parse token: %v", err)
	}

	return &token, nil
}

// RefreshToken refreshes an expired token
func (r *CalendarRepository) RefreshToken(ctx context.Context, token *oauth2.Token, userID int) (*oauth2.Token, error) {
	// Create a new token source using the refresh token
	tokenSource := r.Config.TokenSource(ctx, token)

	// Get a new token
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh token: %v", err)
	}

	// Store the refreshed token in the database
	if err := r.StoreUserToken(userID, newToken); err != nil {
		return nil, fmt.Errorf("unable to store refreshed token: %v", err)
	}

	return newToken, nil
}
