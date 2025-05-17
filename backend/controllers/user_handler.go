package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/johneliud/evently/backend/models"
	"github.com/johneliud/evently/backend/repositories"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	UserRepo *repositories.UserRepository
}

func NewUserHandler(userRepo *repositories.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// SignUp handles user registration
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	var req models.UserSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Println("Invalid request body")
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.ConfirmedPassword) == "" || strings.TrimSpace(req.FirstName) == "" || strings.TrimSpace(req.LastName) == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		log.Println("All fields are required")
		return
	}

	if req.Password != req.ConfirmedPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		log.Println("Passwords do not match")
		return
	}

	// Create user
	id, err := h.UserRepo.CreateUser(req)
	if err != nil {
		if err.Error() == "email already exists" {
			http.Error(w, "Email already exists", http.StatusConflict)
			log.Printf("Email already exists: %v\n", err)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Failed to create user: %v\n", err)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "User created successfully",
	})
	log.Println("User created successfully")
}
