package integration_test

import (
	"database/sql"
	"log/slog"

	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
)

type User struct {
	repository.User
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

func CreateUser(db *sql.DB, username string) (User, error) {
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

func (u User) getAuthToken(secret []byte) (string, error) {
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

func CreateReceiptWithUser(db *sql.DB, total int, description string) (Receipt, error) {
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

	err = tx.QueryRow("INSERT INTO receipts (total, description, user_id) VALUES ($1, $2, $3) RETURNING id", receipt.Total, receipt.Description, receipt.UserID).Scan(&receipt.ID)
	if err != nil {
		slog.Error(err.Error())
		return receipt, err
	}

	user, err := CreateUser(db, "fake_user")
	if err != nil {
		return receipt, err
	}

	receipt.User = user

	return receipt, nil
}
