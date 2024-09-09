package repository

import (
	"database/sql"
	"fmt"

	"github.com/carlqt/ezsplit/.gen/public/model"
	. "github.com/carlqt/ezsplit/.gen/public/table"
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
		Users.ID, Users.Name,
	).FROM(
		Users.
			INNER_JOIN(UserOrders, Users.ID.EQ(UserOrders.UserID)),
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
	var itemIDExpression []Expression
	var userOrders []UserOrder

	// Fetch all the user_orders of the user with ID $user_id
	userOrdersStmt := UserOrders.SELECT(UserOrders.ItemID).WHERE(UserOrders.UserID.EQ(RawInt(userID)))
	err := userOrdersStmt.Query(r.DB, &userOrders)
	if err != nil {
		return 0, fmt.Errorf("failed to get the user orders: %w", err)
	} else if len(userOrders) == 0 {
		return 0, nil
	}

	// Initializing itemIDExpression from the userOrders.ItemID
	for _, userOrder := range userOrders {
		itemID := userOrder.ItemID
		itemIDExpression = append(itemIDExpression, Int(*itemID))
	}

	// Fetch all items and all user_orders associated
	// This is needed to calculate the "shared_price" of the item
	var itemsWithOrders []struct {
		model.Items
		UserOrders []struct {
			model.UserOrders
		}
	}

	itemsStmt := SELECT(
		Items.ID,
		Items.Price,
		UserOrders.AllColumns,
	).FROM(
		Items.INNER_JOIN(
			UserOrders,
			UserOrders.ItemID.EQ(Items.ID),
		),
	).WHERE(Items.ID.IN(itemIDExpression...))

	err = itemsStmt.Query(r.DB, &itemsWithOrders)
	if err != nil {
		return 0, fmt.Errorf("failed to get the user items: %w", err)
	}

	// Calculate the TotalPayables
	for _, item := range itemsWithOrders {
		share := len(item.UserOrders)
		totalPayables = totalPayables + (int(item.Price) / share)
	}

	return totalPayables, err
}

func (r *UserOrdersRepository) FindByUserIDAndItemID(userID string, itemID string) (UserOrder, error) {
	var userOrder UserOrder

	stmt := UserOrders.SELECT(
		UserOrders.AllColumns,
	).FROM(
		UserOrders.INNER_JOIN(Items, UserOrders.ItemID.EQ(Items.ID)).INNER_JOIN(Users, UserOrders.UserID.EQ(Users.ID)),
	).WHERE(
		UserOrders.UserID.EQ(RawInt(userID)).AND(UserOrders.ItemID.EQ(RawInt(itemID))),
	).LIMIT(1)

	err := stmt.Query(r.DB, &userOrder)
	if err != nil {
		return userOrder, fmt.Errorf("failed to find user_order with user_id=%s and item_id=%s: %w", userID, itemID, err)
	}

	return userOrder, nil
}
