package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // The "-" tag ensures the password is not included in JSON output
}

func NewUser(email, password string) *User {
	return &User{
		Email:    email,
		Password: password,
	}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
