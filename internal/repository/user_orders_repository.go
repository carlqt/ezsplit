package repository

import (
	"database/sql"
	"fmt"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
	. "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type UserOrdersRepository struct {
	DB *sql.DB
}

type UserOrder struct {
	model.UserOrders
}

func (r *UserOrdersRepository) Create(userID string, itemID string) error {
  stmt := UserOrders.INSERT(UserOrders.UserID, UserOrders.ItemID).VALUES(userID, itemID)

	_, err := stmt.Exec(r.DB)
	if err != nil {
		return fmt.Errorf("failed to insert user_orders with user_id=%s and item_id=%s: %w", userID, itemID, err)
	}

	return err
}

func (r *UserOrdersRepository) Delete(userID string, itemID string) error {
  stmt := UserOrders.DELETE().WHERE(UserOrders.UserID.EQ(RawInt(userID)).AND(UserOrders.ItemID.EQ(RawInt(itemID))))

	_, err := stmt.Exec(r.DB)
	if err != nil {
    return fmt.Errorf("failed to delete user_order from DB: %w", err)
	}

	return err
}

func (r *UserOrdersRepository) SelectAllUsersFromItem(itemID string) ([]User, error) {
  users := []User{}

  stmt := SELECT(
    Users.ID, Users.Username,
  ).FROM(
    Users.
      INNER_JOIN(UserOrders,Users.ID.EQ(UserOrders.UserID)),
  ).WHERE(
    UserOrders.ItemID.EQ(RawInt(itemID)),
  )

	err := stmt.Query(r.DB, &users)
	if err != nil {
    return nil, fmt.Errorf("failed to fetch associated users of item_id=%s: %w", itemID, err)
	}

	return users, nil
}

func (r *UserOrdersRepository) GetTotalPayables(userID string) (int, error) {
	var totalPayables int
  var itemIDs []Expression
  var items []Item

  // Get all UserOrders where user_id = userID
  userOrders := []UserOrder{}
  userOrdersStmt := UserOrders.SELECT(UserOrders.UserID, UserOrders.ItemID).WHERE(UserOrders.UserID.EQ(RawInt(userID)))
  err := userOrdersStmt.Query(r.DB, &userOrders)
  if err != nil {
    return 0, fmt.Errorf("failed to get the user's total payables")
  }

  // Get All Items from the UserOrders fetched
  for _, userOrder := range userOrders {
    itemIDs = append(itemIDs, Int(*userOrder.ItemID))
  }

  itemsStmt := Items.SELECT(Items.ID, Items.Price).WHERE(Items.ID.IN(itemIDs...))
  err = itemsStmt.Query(r.DB, &items)
  if err != nil {
    return 0, fmt.Errorf("failed to get the user's total payables")
  }

  // Calculate the TotalPayables
  sumOfItemPrice := 0
  for _, item := range items {
    sumOfItemPrice += int(*item.Price)
  }

  totalPayables = sumOfItemPrice / len(userOrders)

	return totalPayables, err
}
