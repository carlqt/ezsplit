package repository

import "database/sql"

type Repository struct {
	UserRepository       *UserRepository
	ReceiptRepository    *ReceiptRepository
	ItemRepository       *ItemRepository
	UserOrdersRepository *UserOrdersRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:       &UserRepository{DB: db},
		ReceiptRepository:    &ReceiptRepository{DB: db},
		ItemRepository:       &ItemRepository{DB: db},
		UserOrdersRepository: &UserOrdersRepository{DB: db},
	}
}
