package models

import "time"

// User represents a user in the system
type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"-"`
	ConfirmedPassword string    `json:"-"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UserSignupRequest represents the data needed for user signup
type UserSignupRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
}

// UserSignInRequest represents the data needed for user signin
type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
