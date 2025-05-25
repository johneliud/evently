package repositories

import (
	"database/sql"
	"log"

	"github.com/johneliud/evently/backend/models"
)

// RSVPRepository handles database operations for RSVPs
type RSVPRepository struct {
	DB *sql.DB
}

func NewRSVPRepository(db *sql.DB) *RSVPRepository {
	return &RSVPRepository{DB: db}
}

// CreateOrUpdateRSVP creates or updates an RSVP
func (r *RSVPRepository) CreateOrUpdateRSVP(eventID, userID int, status string) error {
	// Check if RSVP already exists
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM rsvps 
			WHERE event_id = $1 AND user_id = $2
		)
	`, eventID, userID).Scan(&exists)

	if err != nil {
		log.Printf("Error checking if RSVP exists: %v", err)
		return err
	}

	if exists {
		// Update existing RSVP
		_, err = r.DB.Exec(`
			UPDATE rsvps 
			SET status = $1, updated_at = NOW() 
			WHERE event_id = $2 AND user_id = $3
		`, status, eventID, userID)
	} else {
		// Create new RSVP
		_, err = r.DB.Exec(`
			INSERT INTO rsvps (event_id, user_id, status) 
			VALUES ($1, $2, $3)
		`, eventID, userID, status)
	}

	if err != nil {
		log.Printf("Error creating/updating RSVP: %v", err)
		return err
	}

	return nil
}

// GetRSVPByEventAndUser gets an RSVP by event ID and user ID
func (r *RSVPRepository) GetRSVPByEventAndUser(eventID, userID int) (*models.RSVP, error) {
	var rsvp models.RSVP
	err := r.DB.QueryRow(`
		SELECT id, event_id, user_id, status, created_at, updated_at
		FROM rsvps
		WHERE event_id = $1 AND user_id = $2
	`, eventID, userID).Scan(
		&rsvp.ID,
		&rsvp.EventID,
		&rsvp.UserID,
		&rsvp.Status,
		&rsvp.CreatedAt,
		&rsvp.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No RSVP found
	}

	if err != nil {
		log.Printf("Error getting RSVP: %v", err)
		return nil, err
	}

	return &rsvp, nil
}

// GetRSVPs gets all RSVPs for an event
func (r *RSVPRepository) GetRSVPs(eventID int) ([]models.RSVPWithUser, error) {
	rows, err := r.DB.Query(`
		SELECT r.id, r.event_id, r.user_id, r.status, r.created_at, r.updated_at,
			   u.first_name, u.last_name, u.email
		FROM rsvps r
		JOIN users u ON r.user_id = u.id
		WHERE r.event_id = $1
		ORDER BY r.created_at DESC
	`, eventID)

	if err != nil {
		log.Printf("Error getting RSVPs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rsvps []models.RSVPWithUser
	for rows.Next() {
		var rsvp models.RSVPWithUser
		err := rows.Scan(
			&rsvp.ID,
			&rsvp.EventID,
			&rsvp.UserID,
			&rsvp.Status,
			&rsvp.CreatedAt,
			&rsvp.UpdatedAt,
			&rsvp.FirstName,
			&rsvp.LastName,
			&rsvp.Email,
		)

		if err != nil {
			log.Printf("Error scanning RSVP row: %v", err)
			return nil, err
		}

		rsvps = append(rsvps, rsvp)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating RSVP rows: %v", err)
		return nil, err
	}

	return rsvps, nil
}

// GetRSVPCount gets the count of RSVPs by status for an event
func (r *RSVPRepository) GetRSVPCount(eventID int) (models.RSVPCount, error) {
	var count models.RSVPCount
	count.EventID = eventID

	// Get going count
	err := r.DB.QueryRow(`
		SELECT COUNT(*) FROM rsvps
		WHERE event_id = $1 AND status = 'going'
	`, eventID).Scan(&count.Going)

	if err != nil {
		log.Printf("Error getting going count: %v", err)
		return count, err
	}

	// Get maybe count
	err = r.DB.QueryRow(`
		SELECT COUNT(*) FROM rsvps
		WHERE event_id = $1 AND status = 'maybe'
	`, eventID).Scan(&count.Maybe)

	if err != nil {
		log.Printf("Error getting maybe count: %v", err)
		return count, err
	}

	// Get not going count
	err = r.DB.QueryRow(`
		SELECT COUNT(*) FROM rsvps
		WHERE event_id = $1 AND status = 'not_going'
	`, eventID).Scan(&count.NotGoing)

	if err != nil {
		log.Printf("Error getting not going count: %v", err)
		return count, err
	}

	return count, nil
}

// DeleteRSVP deletes an RSVP
func (r *RSVPRepository) DeleteRSVP(eventID, userID int) error {
	_, err := r.DB.Exec(`
		DELETE FROM rsvps
		WHERE event_id = $1 AND user_id = $2
	`, eventID, userID)

	if err != nil {
		log.Printf("Error deleting RSVP: %v", err)
		return err
	}

	return nil
}
