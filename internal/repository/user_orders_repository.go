package repository

import (
	"database/sql"
	"log/slog"
	"time"
)

type UserOrders struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ItemID    string    `json:"item_id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

type UserOrdersRepository struct {
	DB *sql.DB
}

func (r *UserOrdersRepository) Create(userID string, itemID string) error {
	query := `INSERT INTO user_orders (user_id, item_id) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, userID, itemID)
	if err != nil {
		slog.Error(err.Error())
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
