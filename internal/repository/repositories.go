package repository

import (
	"database/sql"
	"strconv"
)

type Repository struct {
	UserRepository       *UserRepository
	ReceiptRepository    *ReceiptRepository
	ItemRepository       *ItemRepository
	UserOrdersRepository *UserOrdersRepository
	AccountRepository    *AccountRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		UserRepository:       &UserRepository{DB: db},
		ReceiptRepository:    &ReceiptRepository{DB: db},
		ItemRepository:       &ItemRepository{DB: db},
		UserOrdersRepository: &UserOrdersRepository{DB: db},
		AccountRepository:    &AccountRepository{DB: db},
	}
}

// BigInt helps with initializing a repository model that has a BigInt data type in the database
func BigInt[T string | int32](input T) int64 {
	switch v := any(input).(type) {
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	case int32:
		return int64(v)
	default:
		// Not sure the best way to handle this but I'm hoping that the code wouldn't reach here
		return 0
	}
}

// Nullable helps with initializing structs that contains pointer fields
// the pointer fields is because the column in the DB allows for NULL values
func Nullable[T any](input T) *T {
	return &input
}
