package repositories

import (
	"database/sql"
	"fmt"
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
	var id int
	err := r.DB.QueryRow("SELECT id FROM users WHERE email = $1", user.Email).Scan(&id)
	if err == nil && id != 0 {
		log.Println("User already exists with that email:", id, user.Email)
		return id, fmt.Errorf("email already exists")
	}

	fmt.Println("id:", id)
	fmt.Println("userEmail:", user.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return 0, err
	}

	// Insert the user
	err = r.DB.QueryRow(
		"INSERT INTO users (email, password, confirmed_password, first_name, last_name) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Email, string(hashedPassword), string(hashedPassword), user.FirstName, user.LastName,
	).Scan(&id)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return 0, err
	}
	log.Printf("User created with ID: %d", id)
	return id, nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow(
		"SELECT id, email, password, first_name, last_name, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}
	
	return &user, nil
}
