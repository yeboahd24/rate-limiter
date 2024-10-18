package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/yeboahd24/rate-limiter/model"
	"github.com/yeboahd24/rate-limiter/repository"
	"golang.org/x/crypto/bcrypt"
)

// Claims represents the JWT claims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// AuthService handles user authentication
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(email, password string) error {
	// Check if user already exists
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return fmt.Errorf("user with email %s already exists", email)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing user: %v", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	user := model.NewUser(email, string(hashedPassword))
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

// LoginUser logs in an existing user and returns a JWT token
func (s *AuthService) LoginUser(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no user found with email %s", email)
		}
		return "", fmt.Errorf("error fetching user: %v", err)
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	// Generate a JWT token
	claims := Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return tokenString, nil
}

// GetUser retrieves a user by their ID
func (s *AuthService) GetUser(userID int) (*model.User, error) {
	return s.userRepo.GetUserByID(userID)
}
