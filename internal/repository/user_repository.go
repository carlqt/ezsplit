// Memo: The user repository is a bit complicated as it abstracts 2 user structs and also plays around with Struct Enums.
// Abstracting the user is actually the preferred implementation as it allows us to change the underlying model without affecting the public API.

package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/public/model"
	. "github.com/carlqt/ezsplit/.gen/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
)

var ErrWrongCredentials = errors.New("incorrect username or password")

type User struct {
	model.Users
	Account Account
}

type UserRepository struct {
	DB *sql.DB
}

func (u User) IsVerified() bool {
	return u.AccountID != nil
}

// TODO: Introduce a service layer and move this logic there
func (r *UserRepository) CreateWithAccount(username, password string) (User, error) {
	var account Account
	var user User

	tx, err := r.DB.Begin()
	if err != nil {
		return user, fmt.Errorf("failed to start the transaction: %w", err)
	}

	//nolint:errcheck
	defer tx.Rollback()

	account.Username = username
	account.Password, err = hashPassword(password)
	if err != nil {
		slog.Error("failed to create user with account", "error", err.Error())
	}

	accountStmt := Accounts.INSERT(
		Accounts.Username, Accounts.Password,
	).MODEL(account).RETURNING(Accounts.ID, Accounts.Username, Accounts.CreatedAt)

	err = accountStmt.Query(tx, &account)
	if err != nil {
		return user, fmt.Errorf("failed to create the account: %w", err)
	}

	user.Name = username
	user.AccountID = Nullable(account.ID)
	userStmt := Users.INSERT(
		Users.Name, Users.AccountID,
	).MODEL(user).RETURNING(Users.Name, Users.ID, Users.AccountID, Users.CreatedAt)

	err = userStmt.Query(tx, &user)

	if err != nil {
		return user, fmt.Errorf("failed to create the user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return user, fmt.Errorf("failed to commit the transaction: %w", err)
	}

	user.Account = account

	return user, nil
}

func (r *UserRepository) CreateGuest(username string) (User, error) {
	var user User
	user.Name = username

	stmt := Users.INSERT(
		Users.Name,
	).MODEL(user).RETURNING(Users.Name, Users.ID)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("failed to create guest %s: %w", username, err)
	}
	return user, nil
}

func (r *UserRepository) FindByID(id string) (User, error) {
	var user User
	stmt := SELECT(Users.ID, Users.Name, Users.AccountID).FROM(Users.Table).WHERE(Users.ID.EQ(RawInt(id))).LIMIT(1)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("%w: DB Query failed for id=%s", err, id)
	}
	return user, nil
}

func (r *UserRepository) FindByName(name string) (User, error) {
	var user User

	stmt := Users.SELECT(Users.ID, Users.Name, Users.AccountID).WHERE(Users.Name.EQ(String(name))).LIMIT(1)

	err := stmt.Query(r.DB, &user)
	if err != nil {
		return user, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

// FindVerifiedByUsername finds users with an account associated with them
func (r *UserRepository) FindVerifiedByUsername(username, password string) (User, error) {
	var user User

	stmt := SELECT(
		Users.ID, Users.Name, Users.AccountID, Accounts.Password, Accounts.ID, Accounts.Username,
	).FROM(
		Users, Accounts,
	).WHERE(Accounts.Username.EQ(String(username))).LIMIT(1)

	err := stmt.Query(r.DB, &user)

	if err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return user, fmt.Errorf("database error: %w", err)
	}

	if errors.Is(err, qrm.ErrNoRows) || !validatePasswords(user.Account.Password, password) {
		return user, ErrWrongCredentials
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	var returnedUsers []User

	stmt := Users.SELECT(Users.ID, Users.Name, Users.AccountID)

	err := stmt.Query(r.DB, &returnedUsers)
	if err != nil {
		return returnedUsers, fmt.Errorf("failed to get users: %w", err)
	}

	return returnedUsers, nil
}
