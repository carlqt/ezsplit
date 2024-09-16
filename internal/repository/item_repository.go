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
	stmt := Items.SELECT(Items.ID, Items.Name, Items.Price, Items.ReceiptID, Items.CreatedAt).WHERE(Items.ID.EQ(RawInt(id))).LIMIT(1)

	err := stmt.Query(i.DB, &item)
	if err != nil {
		return item, fmt.Errorf("failed to find item with id=%s: %w", id, err)
	}

	return item, nil
}

func (i *ItemRepository) GetMyOrders(userId, receiptID string) ([]Item, error) {
	var items []Item

	condition := Bool(true)

	// If receiptID is given, filter the query to the receipt
	if receiptID != "" {
		condition = condition.AND(Items.ReceiptID.EQ(RawInt(receiptID)))
	}

	stmt := SELECT(
		Items.ID, Items.Name, Items.Price, Items.ReceiptID,
	).FROM(
		UserOrders.INNER_JOIN(Items, UserOrders.ItemID.EQ(Items.ID)).INNER_JOIN(Users, UserOrders.UserID.EQ(Users.ID)),
	).WHERE(condition.AND(UserOrders.UserID.EQ(RawInt(userId))))

	err := stmt.Query(i.DB, &items)
	if err != nil {
		return items, fmt.Errorf("failed to get the orders: %w", err)
	}

	return items, nil
}

func (i *ItemRepository) DeleteFromReceipt(userID, itemID string) (Item, error) {
	var item Item

  stmt := Items.DELETE().USING(Receipts).WHERE(
    Items.ID.EQ(
      RawInt(itemID),
    ).AND(
      Items.ReceiptID.EQ(Receipts.ID),
    ).AND(
      Receipts.UserID.EQ(RawInt(userID)),
    ),
  ).RETURNING(Items.AllColumns)

  err := stmt.Query(i.DB, &item)
  if err != nil {
    return item, fmt.Errorf("failed to delete item: %w", err)
  }

  return item, nil
}
