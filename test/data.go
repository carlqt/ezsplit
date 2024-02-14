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

func CreateUser(db *sql.DB, username string) (User, error) {
	user := User{
		repository.User{
			Username: username,
		},
	}

	sql := "INSERT INTO users (username) VALUES ($1) RETURNING id"
	err := db.QueryRow(sql, username).Scan(&user.ID)
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
}

func CreateReceiptWithUser(db *sql.DB, userID string, total int, description string) (Receipt, error) {

	receipt := Receipt{
		repository.Receipt{
			UserID:      userID,
			Total:       total,
			Description: description,
		},
	}

	err := r.DB.QueryRow("INSERT INTO receipts (total, description, user_id) VALUES ($1, $2, $3) RETURNING id", receipt.Total, receipt.Description, receipt.UserID).Scan(&receipt.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}
