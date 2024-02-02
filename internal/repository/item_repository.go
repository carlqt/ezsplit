package repository

import (
	"database/sql"
	"time"
)

type ItemRepository struct {
	DB *sql.DB
}

type Item struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Price     int       `db:"price"`
	ReceiptID string    `db:"receipt_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (i *ItemRepository) Create(item *Item) error {
	err := i.DB.QueryRow("INSERT INTO items (name, price, receipt_id) VALUES ($1, $2, $3) RETURNING id", item.Name, item.Price, item.ReceiptID).Scan(&item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (i *ItemRepository) SelectAllForReceipt(receiptID string) ([]*Item, error) {
	rows, err := i.DB.Query("SELECT id, name, price, created_at FROM items WHERE receipt_id = $1", receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
