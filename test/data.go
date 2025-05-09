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

func CreateVerifiedUser(db DbWriter, username string) (User, error) {
	var user User
	var account repository.Account

	fakePassword := "password"
	hashedPassword, _ := auth.HashPassword(fakePassword)

	sql := "INSERT INTO accounts (username, password) VALUES ($1, $2) RETURNING id"
	err := db.QueryRow(sql, username, hashedPassword).Scan(&account.ID)
	if err != nil {
		return user, err
	}

	user.Name = username
	user.AccountID = repository.Nullable(account.ID)

	sql = "INSERT INTO users (name, account_id) VALUES ($1, $2) RETURNING id, name, account_id"
	err = db.QueryRow(sql, username, user.AccountID).Scan(&user.ID, &user.Name, &user.AccountID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u User) GetAuthToken(secret []byte) (string, error) {
	userClaim := auth.NewUserClaim(u.ID, u.Name, u.IsVerified())
	accessToken, err := auth.CreateAndSignToken(userClaim, secret)
	if err != nil {
		return accessToken, nil
	}

	return accessToken, err
}

type ReceiptWithUser struct {
	repository.Receipt
	User User
}

func CreateReceiptWithUser(db *sql.DB, total int, description string) (ReceiptWithUser, error) {
	receipt := ReceiptWithUser{}
	receipt.Receipt = repository.Receipt{}
	receipt.Description = description
	receiptTotal := int32(total)
	receipt.Total = &receiptTotal

	tx, err := db.Begin()
	if err != nil {
		return receipt, err
	}

	username := fmt.Sprintf("fake_user+%d", rand.IntN(100))
	user, err := CreateVerifiedUser(tx, username)
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

	receipt.User = user

	return receipt, nil
}
