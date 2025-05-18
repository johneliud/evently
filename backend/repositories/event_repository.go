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

