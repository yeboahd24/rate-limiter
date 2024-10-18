package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yeboahd24/rate-limiter/model"
	"github.com/yeboahd24/rate-limiter/repository"
	services "github.com/yeboahd24/rate-limiter/service"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	return &AuthHandler{authService: authService}
}

// func (h *AuthHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	// Extract the user ID from the JWT token
// 	claims, ok := r.Context().Value("user").(services.Claims)
// 	log.Printf("Claims: %v", claims)
// 	if !ok {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	// Fetch the user from the database
// 	user, err := h.authService.GetUser(claims.UserID)
// 	if err != nil {
// 		http.Error(w, "Error fetching user", http.StatusInternalServerError)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(user)
// }

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add logging
	log.Printf("Attempting to log in user with email: %s", user.Email)

	token, err := h.authService.LoginUser(user.Email, user.Password)
	if err != nil {
		log.Printf("Error logging in user: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add logging
	log.Printf("Attempting to register user with email: %s", user.Email)

	err = h.authService.RegisterUser(user.Email, user.Password)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, fmt.Sprintf("Error registering user: %v", err), http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user claims from the JWT token
	claims, ok := r.Context().Value("user").(*model.Claims)

	log.Printf("Retrieving claims from context. Type: %T, Value: %+v", r.Context().Value("user"), r.Context().Value("user"))
	if !ok {
		log.Printf("Failed to extract claims from context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Claims: %v", claims)

	// Fetch the user from the database
	user, err := h.authService.GetUser(claims.UserID)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
