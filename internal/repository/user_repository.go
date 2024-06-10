package repository

import (
	"database/sql"
	"log/slog"
	"time"
)

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(username string, password string) (User, error) {
	user := User{Username: username}

	err := r.DB.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, password).Scan(&user.ID)
	if err != nil {
		slog.Error(err.Error())
		return user, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	user := &User{}
	err := r.DB.QueryRow("SELECT id, username FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]*User, error) {
	rows, err := r.DB.Query("SELECT id, username FROM users")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
