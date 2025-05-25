package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/johneliud/evently/backend/models"
	"github.com/johneliud/evently/backend/repositories"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
		log.Printf("Invalid request body: %v\n", err)
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
	id, err := h.UserRepo.CreateUser(models.UserRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if err.Error() == "email already exists" {
			http.Error(w, "Email already exists", http.StatusConflict)
			log.Printf("Email already exists: %v\n", id)
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

// SignIn handles user authentication
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	var req models.UserSignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Invalid request body: %v\n", err)
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		log.Println("Email and password are required")
		return
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("Invalid email or password: %v\n", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("Invalid email or password: %v\n", err)
		return
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Failed to generate token: %v\n", err)
		return
	}

	// Return success response with token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":   tokenString,
		"user_id": user.ID,
		"message": "Login successful",
	})
	log.Println("Login successful")
}

// GoogleAuthURL returns the URL for Google OAuth
func (h *UserHandler) GoogleAuthURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Generate a random state token to prevent request forgery
	state := fmt.Sprintf("auth-%d", time.Now().UnixNano())

	// Get client ID and secret from environment variables
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Println("WARNING: GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET environment variable is not set")
		http.Error(w, "OAuth configuration error", http.StatusInternalServerError)
		return
	}

	// Create OAuth config
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:9000/api/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Store state in a cookie for verification
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   3600, // 1 hour
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Get the authorization URL
	authURL := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Return the authorization URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"auth_url": authURL,
	})
}

// GoogleCallback handles the OAuth callback from Google
func (h *UserHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println("Method not allowed")
		return
	}

	// Get the state and code from the query parameters
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")

	if state == "" || code == "" {
		log.Println("Missing state or code parameter")
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Verify state to prevent CSRF
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Error retrieving state cookie: %v", err)
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	log.Printf("State from cookie: %s, State from query: %s", stateCookie.Value, state)

	if stateCookie.Value != state {
		log.Printf("State mismatch: cookie=%s, query=%s", stateCookie.Value, state)
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Get client ID and secret from environment variables
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Println("Google OAuth credentials not set")
		http.Error(w, "OAuth configuration error", http.StatusInternalServerError)
		return
	}

	// Create OAuth config
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:9000/api/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Exchange the authorization code for a token
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		log.Printf("Failed to exchange token: %v\n", err)
		return
	}

	// Get user info from Google
	client := config.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		log.Printf("Failed to get user info: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Parse user info
	var userInfo struct {
		Email     string `json:"email"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Picture   string `json:"picture"`
		Sub       string `json:"sub"` // Google's unique user ID
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		log.Printf("Failed to parse user info: %v\n", err)
		return
	}

	// Check if user exists
	user, err := h.UserRepo.GetUserByEmail(userInfo.Email)
	if err != nil {
		// User doesn't exist, create a new one
		// Generate a random password for Google users
		randomPassword := fmt.Sprintf("google_%d", time.Now().UnixNano())

		// Create user
		id, err := h.UserRepo.CreateUser(models.UserRequest{
			Email:     userInfo.Email,
			Password:  randomPassword,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
		})
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			log.Printf("Failed to create user: %v\n", err)
			return
		}

		// Get the newly created user
		user, err = h.UserRepo.GetUserByID(id)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			log.Printf("Failed to get user: %v\n", err)
			return
		}
	}

	// Create JWT token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with a secret key
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Failed to generate token: %v\n", err)
		return
	}

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("http://localhost:5173/auth/callback?token=%s&user_id=%d", tokenString, user.ID)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
