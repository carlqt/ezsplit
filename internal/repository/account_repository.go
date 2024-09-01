package repository

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/public/model"
)

type Account struct {
	model.Accounts
}

type AccountRepository struct {
	DB *sql.DB
}
