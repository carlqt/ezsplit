package integration_test

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand/v2"

	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
)

type User struct {
	repository.User
}

type DbWriter interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

func TruncateAllTables(db *sql.DB) {
	query := `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname =current_schema()) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`

	slog.Debug("Clearing data")

	_, err := db.Exec(query)
	if err != nil {
		slog.Error(err.Error())
	}
}

func CreateUser(db DbWriter, username string) (User, error) {
	fakePassword := "password"
	hashedPassword, _ := auth.HashPassword(fakePassword)

	user := User{
		repository.User{
			Username: username,
		},
	}

	sql := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(sql, username, hashedPassword).Scan(&user.ID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u User) GetAuthToken(secret []byte) (string, error) {
	userClaim := auth.UserClaim{
		ID:       u.ID,
		Username: u.Username,
	}
	accessToken, err := auth.CreateAndSignToken(userClaim, secret)
	if err != nil {
		return accessToken, nil
	}

	return accessToken, err
}

type Receipt struct {
	repository.Receipt
	User User
}

func CreateRandomUser(db DbWriter) string {
	var id string

	stmt := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"

	username := fmt.Sprintf("fake_user+%d", rand.IntN(100))
	err := db.QueryRow(stmt, username, "password").Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

func CreateReceiptWithUser(db *sql.DB, total int, description string) (Receipt, error) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	userID := CreateRandomUser(tx)
	receipt := Receipt{
		Receipt: repository.Receipt{
			Description: description,
			UserID:      userID,
			Total:       total,
		},
	}

	stmt := "INSERT INTO receipts (user_id, total, description) VALUES ($1, $2, $3) RETURNING ID"
	err = tx.QueryRow(stmt, receipt.UserID, receipt.Total, receipt.Description).Scan(&receipt.ID)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	return receipt, nil
}

func CreateReceiptWithUserxxxx(db *sql.DB, total int, description string) (Receipt, error) {
	receipt := Receipt{
		repository.Receipt{
			Total:       total,
			Description: description,
		},
		User{},
	}

	tx, err := db.Begin()
	if err != nil {
		return receipt, err
	}

	username := fmt.Sprintf("fake_user+%d", rand.IntN(100))
	user, err := CreateUser(tx, username)
	if err != nil {
		tx.Rollback()
		return receipt, err
	}

	err = tx.QueryRow("INSERT INTO receipts (total, description, user_id) VALUES ($1, $2, $3) RETURNING id", receipt.Total, receipt.Description, user.ID).Scan(&receipt.ID)
	if err != nil {
		slog.Error(err.Error())
		tx.Rollback()
		return receipt, err
	}

	if err = tx.Commit(); err != nil {
		return receipt, err
	}

	receipt.UserID = user.ID
	receipt.User = user

	return receipt, nil
}
