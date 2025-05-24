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

// CreateOrUpdateRSVP creates or updates a user's RSVP for an event
func (r *RSVPRepository) CreateOrUpdateRSVP(eventID, userID int, status string) error {
    _, err := r.DB.Exec(`
        INSERT INTO event_rsvps (event_id, user_id, status)
        VALUES ($1, $2, $3)
        ON CONFLICT (event_id, user_id) 
        DO UPDATE SET status = $3, updated_at = NOW()
    `, eventID, userID, status)

    if err != nil {
        log.Printf("Error creating/updating RSVP: %v", err)
        return err
    }

    return nil
}

// GetRSVPByEventAndUser gets a user's RSVP for an event
func (r *RSVPRepository) GetRSVPByEventAndUser(eventID, userID int) (*models.RSVP, error) {
    var rsvp models.RSVP
    err := r.DB.QueryRow(`
        SELECT id, event_id, user_id, status, created_at, updated_at
        FROM event_rsvps
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

// DeleteRSVP deletes a user's RSVP for an event
func (r *RSVPRepository) DeleteRSVP(eventID, userID int) error {
    _, err := r.DB.Exec("DELETE FROM event_rsvps WHERE event_id = $1 AND user_id = $2", eventID, userID)
    if err != nil {
        log.Printf("Error deleting RSVP: %v", err)
        return err
    }
    return nil
}

// GetRSVPCountByEvent gets the count of RSVPs by status for an event
func (r *RSVPRepository) GetRSVPCountByEvent(eventID int) (*models.RSVPCount, error) {
    var count models.RSVPCount
    
    err := r.DB.QueryRow(`
        SELECT 
            COUNT(CASE WHEN status = 'going' THEN 1 END) as going,
            COUNT(CASE WHEN status = 'maybe' THEN 1 END) as maybe,
            COUNT(CASE WHEN status = 'not_going' THEN 1 END) as not_going
        FROM event_rsvps
        WHERE event_id = $1
    `, eventID).Scan(
        &count.Going,
        &count.Maybe,
        &count.NotGoing,
    )

    if err != nil {
        log.Printf("Error getting RSVP count: %v", err)
        return nil, err
    }

    return &count, nil
}

// GetRSVPsByEvent gets all RSVPs for an event with user information
func (r *RSVPRepository) GetRSVPsByEvent(eventID int) ([]models.RSVPWithUser, error) {
    rows, err := r.DB.Query(`
        SELECT r.id, r.event_id, r.user_id, r.status, r.created_at, r.updated_at,
               u.first_name, u.last_name
        FROM event_rsvps r
        JOIN users u ON r.user_id = u.id
        WHERE r.event_id = $1
        ORDER BY r.created_at DESC
    `, eventID)
    
    if err != nil {
        log.Printf("Error getting RSVPs: %v", err)
        return nil, err
    }
    defer rows.Close()

    rsvps := []models.RSVPWithUser{}
    
    for rows.Next() {
        var rsvp models.RSVPWithUser
        if err := rows.Scan(
            &rsvp.ID,
            &rsvp.EventID,
            &rsvp.UserID,
            &rsvp.Status,
            &rsvp.CreatedAt,
            &rsvp.UpdatedAt,
            &rsvp.FirstName,
            &rsvp.LastName,
        ); err != nil {
            log.Printf("Error scanning RSVP row: %v", err)
            return nil, err
        }
        rsvps = append(rsvps, rsvp)
    }

    return rsvps, nil
}