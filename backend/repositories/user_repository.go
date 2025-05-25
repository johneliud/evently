package repositories

import (
	"database/sql"
	"log"

	"github.com/johneliud/evently/backend/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository handles database operations for users
type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// GetUserByID gets a user by ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(`
		SELECT id, first_name, last_name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail gets a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(`
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user models.UserRequest) (int, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return 0, err
	}

	var id int
	err = r.DB.QueryRow(`
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, user.FirstName, user.LastName, user.Email, string(hashedPassword)).Scan(&id)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0, err
	}

	return id, nil
}
