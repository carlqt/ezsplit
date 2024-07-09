package repository

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
	. "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
)

type ItemRepository struct {
	DB *sql.DB
}

type Item struct {
	model.Items
}

func (i *ItemRepository) Create(item *Item) error {
  stmt := Items.INSERT(Items.Name, Items.Price, Items.ReceiptID).VALUES(item.Name, item.Price, item.ReceiptID).RETURNING(Items.Name, Items.Price, Items.ReceiptID, Items.ID)

	err := stmt.Query(i.DB, item)
	if err != nil {
    return fmt.Errorf("failed to create item in db: %w", err)
	}

	return nil
}

func (i *ItemRepository) SelectAllForReceipt(receiptID string) ([]*Item, error) {
	rows, err := i.DB.Query("SELECT id, name, price, created_at FROM items WHERE receipt_id = $1", receiptID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item := &Item{}
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.CreatedAt)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (i *ItemRepository) FindByID(id string) (*Item, error) {
	item := &Item{}
	err := i.DB.QueryRow("SELECT id, name, price, receipt_id, created_at FROM items WHERE id = $1", id).Scan(&item.ID, &item.Name, &item.Price, &item.ReceiptID, &item.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	return item, nil
}
