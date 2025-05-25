package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/johneliud/evently/backend/models"
)

// EventRepository handles database operations for events
type EventRepository struct {
	DB *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{DB: db}
}

// CreateEvent creates a new event in the database
func (r *EventRepository) CreateEvent(event models.EventRequest, userID int) (int, error) {
	var id int
	err := r.DB.QueryRow(
		"INSERT INTO events (title, description, date, location, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		event.Title, event.Description, event.Date, event.Location, userID,
	).Scan(&id)

	if err != nil {
		log.Printf("Error creating event: %v", err)
		return 0, err
	}

	return id, nil
}

// GetEventsByUserID retrieves all events for a specific user
func (r *EventRepository) GetEventsByUserID(userID int) ([]models.Event, error) {
	rows, err := r.DB.Query(
		"SELECT id, title, description, date, location, user_id, created_at, updated_at FROM events WHERE user_id = $1 ORDER BY date",
		userID,
	)
	if err != nil {
		log.Printf("Error getting events: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Location,
			&event.UserID,
			&event.CreatedAt,
			&event.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning event row: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

// GetUpcomingEvents retrieves all upcoming events
func (r *EventRepository) GetUpcomingEvents() ([]models.Event, error) {
	rows, err := r.DB.Query(`
		SELECT e.id, e.title, e.description, e.date, e.location, e.user_id, e.created_at, e.updated_at,
			   u.first_name, u.last_name
		FROM events e
		JOIN users u ON e.user_id = u.id
		WHERE e.date > NOW()
		ORDER BY e.date ASC
		LIMIT 20
	`)
	if err != nil {
		log.Printf("Error getting upcoming events: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := []models.Event{}

	for rows.Next() {
		var event models.EventWithOrganizer
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Location,
			&event.UserID,
			&event.CreatedAt,
			&event.UpdatedAt,
			&event.OrganizerFirstName,
			&event.OrganizerLastName,
		); err != nil {
			log.Printf("Error scanning event row: %v", err)
			return nil, err
		}
		events = append(events, event.Event)
	}

	return events, nil
}

// GetEventByID retrieves a single event by ID with organizer information
func (r *EventRepository) GetEventByID(id int) (*models.EventWithOrganizer, error) {
	var event models.EventWithOrganizer
	err := r.DB.QueryRow(`
		SELECT e.id, e.title, e.description, e.date, e.location, e.user_id, e.created_at, e.updated_at,
			   u.first_name, u.last_name
		FROM events e
		JOIN users u ON e.user_id = u.id
		WHERE e.id = $1
	`, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.Location,
		&event.UserID,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.OrganizerFirstName,
		&event.OrganizerLastName,
	)

	if err != nil {
		log.Printf("Error getting event by ID: %v", err)
		return nil, err
	}

	return &event, nil
}

// DeleteEvent deletes an event by ID
func (r *EventRepository) DeleteEvent(eventID int) error {
	_, err := r.DB.Exec("DELETE FROM events WHERE id = $1", eventID)
	if err != nil {
		log.Printf("Error deleting event: %v", err)
		return err
	}
	return nil
}

// UpdateEvent updates an existing event
func (r *EventRepository) UpdateEvent(eventID int, event models.EventRequest) error {
	_, err := r.DB.Exec(
		"UPDATE events SET title = $1, description = $2, date = $3, location = $4, updated_at = NOW() WHERE id = $5",
		event.Title, event.Description, event.Date, event.Location, eventID,
	)
	if err != nil {
		log.Printf("Error updating event: %v", err)
		return err
	}
	return nil
}

// SearchEvents searches for events based on title, location, and date range
func (r *EventRepository) SearchEvents(query string, location string, startDate, endDate *time.Time) ([]models.EventWithOrganizer, error) {
	// Build the query dynamically based on provided filters
	queryBuilder := `
		SELECT e.id, e.title, e.description, e.date, e.location, e.user_id, e.created_at, e.updated_at,
			   u.first_name, u.last_name
		FROM events e
		JOIN users u ON e.user_id = u.id
		WHERE 1=1
	`
	var args []interface{}
	argPosition := 1

	// Add title search if query is provided
	if query != "" {
		queryBuilder += fmt.Sprintf(" AND e.title ILIKE $%d", argPosition)
		args = append(args, "%"+query+"%")
		argPosition++
	}

	// Add location filter if provided
	if location != "" {
		queryBuilder += fmt.Sprintf(" AND e.location ILIKE $%d", argPosition)
		args = append(args, "%"+location+"%")
		argPosition++
	}

	// Add date range filters if provided
	if startDate != nil {
		queryBuilder += fmt.Sprintf(" AND e.date >= $%d", argPosition)
		args = append(args, startDate)
		argPosition++
	}

	if endDate != nil {
		queryBuilder += fmt.Sprintf(" AND e.date <= $%d", argPosition)
		args = append(args, endDate)
		argPosition++
	}

	// Only show future events by default if no date filters are provided
	if startDate == nil && endDate == nil {
		queryBuilder += " AND e.date >= NOW()"
	}

	// Order by date
	queryBuilder += " ORDER BY e.date ASC LIMIT 100"

	// Execute the query
	rows, err := r.DB.Query(queryBuilder, args...)
	if err != nil {
		log.Printf("Error searching events: %v", err)
		return nil, err
	}
	defer rows.Close()

	events := []models.EventWithOrganizer{}

	for rows.Next() {
		var event models.EventWithOrganizer
		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Description,
			&event.Date,
			&event.Location,
			&event.UserID,
			&event.CreatedAt,
			&event.UpdatedAt,
			&event.OrganizerFirstName,
			&event.OrganizerLastName,
		); err != nil {
			log.Printf("Error scanning event row: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}
