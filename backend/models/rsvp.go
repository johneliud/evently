package models

import "time"

// RSVP represents an RSVP in the system
type RSVP struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"` // going, maybe, not_going
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RSVPWithUser extends RSVP with user information
type RSVPWithUser struct {
	RSVP
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// RSVPCount represents the count of RSVPs by status
type RSVPCount struct {
	EventID  int `json:"event_id"`
	Going    int `json:"going"`
	Maybe    int `json:"maybe"`
	NotGoing int `json:"not_going"`
}

// RSVPRequest represents the data needed to create or update an RSVP
type RSVPRequest struct {
	Status string `json:"status"` // going, maybe, not_going
}
