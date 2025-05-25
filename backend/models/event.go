package models

import "time"

// Event represents an event in the system
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventWithOrganizer extends Event with organizer information
type EventWithOrganizer struct {
	Event
	OrganizerFirstName string `json:"organizer_first_name"`
	OrganizerLastName  string `json:"organizer_last_name"`
}

// EventRequest represents the data needed to create or update an event
type EventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}
