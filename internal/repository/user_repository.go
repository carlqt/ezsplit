package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/public/model"
	. "github.com/carlqt/ezsplit/.gen/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type userState string

const (
	authenticated userState = "authenticated"
	guest         userState = "guest"
)

type User struct {
	model.Users
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(username string, password string) (User, error) {
	user := User{}

	stmt := Users.INSERT(
		Users.Username, Users.Password, Users.State,
	).VALUES(
		username, password, authenticated,
	).RETURNING(Users.Username, Users.ID)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("failed to create user with username %s: %w", username, err)
	}
	return user, nil
}

func (r *UserRepository) CreateGuest(username string) (User, error) {
	user := User{}

	stmt := Users.INSERT(
		Users.Username, Users.Password, Users.State,
	).VALUES(
		username, "", guest,
	).RETURNING(Users.Username, Users.ID)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("failed to create guest %s: %w", username, err)
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (User, error) {
	user := User{}
	stmt := SELECT(Users.ID, Users.Username).FROM(Users.Table).WHERE(Users.ID.EQ(RawInt(id)))

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("%w: DB Query failed for id=%s", err, id)
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (User, error) {
	user := User{}

	stmt := Users.SELECT(Users.ID, Users.Username, Users.Password).WHERE(Users.Username.EQ(String(username)))

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	users := []User{}
	stmt := Users.SELECT(Users.ID, Users.Username)

	err := stmt.Query(r.DB, &users)
	if err != nil {
		return users, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}
