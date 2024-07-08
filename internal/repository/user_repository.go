package repository

import (
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
	. "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
	. "github.com/go-jet/jet/v2/postgres"
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
