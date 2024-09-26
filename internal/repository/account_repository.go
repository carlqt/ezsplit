package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/carlqt/ezsplit/.gen/public/model"
)

type Account struct {
	model.Accounts
}

type AccountRepository struct {
	DB *sql.DB
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot generate hash from string: %w", err)
	}

	return string(hash), nil
}

func validatePasswords(password string, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword)); err != nil {
		return false
	}

	return true
}
