package repository

import (
	"database/sql"
	"fmt"

	"github.com/carlqt/ezsplit/.gen/public/model"
	. "github.com/carlqt/ezsplit/.gen/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type ItemRepository struct {
	DB *sql.DB
}

type Item struct {
	model.Items
}

func (i *ItemRepository) Create(item *Item) error {
	stmt := Items.INSERT(Items.Name, Items.Price, Items.ReceiptID).MODEL(item).RETURNING(Items.Name, Items.Price, Items.ReceiptID, Items.ID)

	err := stmt.Query(i.DB, item)
	if err != nil {
		return fmt.Errorf("failed to create item in db: %w", err)
	}

	return nil
}

func (i *ItemRepository) SelectAllForReceipt(receiptID string) ([]Item, error) {
	items := []Item{}
	stmt := Items.SELECT(Items.ID, Items.Name, Items.Price, Items.CreatedAt).WHERE(Items.ReceiptID.EQ(RawInt(receiptID)))

	err := stmt.Query(i.DB, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items")
	}

	return items, nil
}

func (i *ItemRepository) FindByID(id string) (Item, error) {
	item := Item{}
	stmt := Items.SELECT(Items.ID, Items.Name, Items.Price, Items.ReceiptID, Items.CreatedAt).WHERE(Items.ID.EQ(RawInt(id)))

	err := stmt.Query(i.DB, &item)
	if err != nil {
		return item, fmt.Errorf("failed to find item with id=%s: %w", id, err)
	}

	return item, nil
}
