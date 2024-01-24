package repository

import "database/sql"

type Repository struct {
	UserRepository *UserRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository: &UserRepository{DB: db},
	}
}
