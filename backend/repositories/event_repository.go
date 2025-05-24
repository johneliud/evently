package repositories

import (
	"database/sql"
	"log"

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
