package repository

import (
	"database/sql"

	"github.com/yeboahd24/rate-limiter/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)
	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow("SELECT id, email, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow("SELECT id, email FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
