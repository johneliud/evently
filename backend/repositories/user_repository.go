package repositories

import (
	"database/sql"
	"errors"
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

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user models.UserSignupRequest) (int, error) {
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.Email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if email exists: %v", err)
		return 0, err
	}
	if exists {
		log.Println("Email already exists")
		return 0, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return 0, err
	}

	// Insert the user
	var id int
	err = r.DB.QueryRow(
		"INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Email, string(hashedPassword), user.FirstName, user.LastName,
	).Scan(&id)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
	}
	log.Printf("User created with ID: %d", id)
	return id, err
}
