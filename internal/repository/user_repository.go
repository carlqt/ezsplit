package repository

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
	. "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
)

type User struct {
	model.Users
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(username string, password string) (User, error) {
	user := User{}

	stmt := Users.INSERT(Users.Username, Users.Password).VALUES(username, password).RETURNING(Users.Username, Users.ID, Users.Username)

  err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("%w | failed to insert username %s", err, username)
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*User, error) {
	user := &User{}
	err := r.DB.QueryRow("SELECT id, username FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, fmt.Errorf("%w: DB Query failed for id=%s", err, id)
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (User, error) {
	user := User{}
	err := r.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		slog.Error(err.Error())
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]*User, error) {
	rows, err := r.DB.Query("SELECT id, username FROM users")
	if err != nil {
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
