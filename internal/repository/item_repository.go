package repository

import (
	"database/sql"
	"time"
)

type ItemRepository struct {
	DB *sql.DB
}

type Item struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Price     int       `db:"price"`
	ReceiptID int       `db:"receipt_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (i *ItemRepository) Create(item *Item) error {
	err := i.DB.QueryRow("INSERT INTO items (name, price, receipt_id) VALUES ($1, $2, $3) RETURNING id", item.Name, item.Price, item.ReceiptID).Scan(&item.ID)
	if err != nil {
		return err
	}
	return nil
}
