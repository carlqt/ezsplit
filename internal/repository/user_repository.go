// Memo: The user repository is a bit complicated as it abstracts 2 user structs and also plays around with Struct Enums.
// Abstracting the user is actually the preferred implementation as it allows us to change the underlying model without affecting the public API.

package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/public/model"
	. "github.com/carlqt/ezsplit/.gen/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type UserState struct {
	value string
}

var (
	Guest    = UserState{"guest"}
	Verified = UserState{"verified"}
)

func (u UserState) String() string {
	return u.value
}

// User is a public struct that will be returned or received by this package
type User struct {
	model.Users
	State UserState
}

// StateVal is a helper method to get the string value of the UserState
func (u User) StateVal() string {
	return u.State.value
}

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(username string, password string) (User, error) {
	user := model.Users{}
	user.Username = username
	user.Password = password
	user.State = Verified.String()

	stmt := Users.INSERT(
		Users.Username, Users.Password, Users.State,
	).MODEL(user).RETURNING(Users.Username, Users.ID)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return convertUserModelToUser(user), fmt.Errorf("failed to create user with username %s: %w", username, err)
	}
	return convertUserModelToUser(user), nil
}

func (r *UserRepository) CreateGuest(username string) (User, error) {
	user := model.Users{}
	user.Username = username
	user.State = Guest.String()

	stmt := Users.INSERT(
		Users.Username, Users.Password, Users.State,
	).MODEL(user).RETURNING(Users.Username, Users.ID)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return convertUserModelToUser(user), fmt.Errorf("failed to create guest %s: %w", username, err)
	}
	return convertUserModelToUser(user), nil
}

func (r *UserRepository) FindByID(id string) (User, error) {
	user := model.Users{}
	stmt := SELECT(Users.ID, Users.Username).FROM(Users.Table).WHERE(Users.ID.EQ(RawInt(id)))

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return convertUserModelToUser(user), fmt.Errorf("%w: DB Query failed for id=%s", err, id)
	}
	return convertUserModelToUser(user), nil
}

func (r *UserRepository) FindByUsername(username string) (User, error) {
	user := model.Users{}

	stmt := Users.SELECT(Users.ID, Users.Username, Users.Password).WHERE(Users.Username.EQ(String(username)))

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return convertUserModelToUser(user), fmt.Errorf("failed to find user: %w", err)
	}

	return convertUserModelToUser(user), nil
}

func (r *UserRepository) FindVerifiedByUsername(username string) (User, error) {
	user := model.Users{}

	stmt := Users.SELECT(
		Users.ID, Users.Username, Users.Password,
	).WHERE(Users.Username.EQ(String(username)).AND(Users.State.EQ(String(Verified.String()))))

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return convertUserModelToUser(user), fmt.Errorf("failed to find user: %w", err)
	}

	return convertUserModelToUser(user), nil
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	var returnedUsers []User

	users := []model.Users{}
	stmt := Users.SELECT(Users.ID, Users.Username)

	err := stmt.Query(r.DB, &users)
	if err != nil {
		return returnedUsers, fmt.Errorf("failed to get users: %w", err)
	}

	for _, u := range users {
		returnedUsers = append(returnedUsers, convertUserModelToUser(u))
	}

	return returnedUsers, nil
}

func convertUserModelToUser(user model.Users) User {
	return User{
		Users: user,
		State: UserState{user.State},
	}
}
