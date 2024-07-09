package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type UserOrdersRepository struct {
	DB *sql.DB
}

func (r *UserOrdersRepository) Create(userID string, itemID string) error {
	query := `INSERT INTO user_orders (user_id, item_id) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, userID, itemID)
	if err != nil {
		return fmt.Errorf("failed to insert user_orders with user_id=%s and item_id=%s: %w", userID, itemID, err)
	}

	return err
}

func (r *UserOrdersRepository) Delete(userID string, itemID string) error {
	query := `DELETE FROM user_orders WHERE user_id = $1 AND item_id = $2`
	sqlResult, err := r.DB.Exec(query, userID, itemID)
	if err != nil {
		slog.Error(err.Error())
	}

	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		slog.Error(err.Error())
	} else if rowsAffected == 0 {
		slog.Warn("No rows affected")
	}

	return err
}

func (r *UserOrdersRepository) SelectAllUsersFromItem(itemID string) ([]*User, error) {
	query := "select users.id, users.username from users inner join user_orders on users.id = user_orders.user_id where user_orders.item_id = $1"

	rows, err := r.DB.Query(query, itemID)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserOrdersRepository) GetTotalPayables(userID string) (int, error) {
	var totalPayables int

	query := `SELECT div(items.price, item_count.shared_by_count) as total_payables
	FROM items
	JOIN user_orders as uo ON items.id = uo.item_id
	JOIN (
		SELECT count(*) as shared_by_count, item_id FROM user_orders GROUP BY item_id
	) AS item_count on item_count.item_id = uo.item_id
	WHERE uo.user_id = $1`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	for rows.Next() {
		var p int

		err := rows.Scan(&p)
		if err != nil {
			slog.Error(err.Error())
			return 0, err
		}

		totalPayables += p
	}

	return totalPayables, err
}
