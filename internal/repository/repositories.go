package repository

import "database/sql"

type Repository struct {
	UserRepository    *UserRepository
	ReceiptRepository *ReceiptRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:    &UserRepository{DB: db},
		ReceiptRepository: &ReceiptRepository{DB: db},
	}
}
