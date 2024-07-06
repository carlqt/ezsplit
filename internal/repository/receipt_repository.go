package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	_ "github.com/lib/pq"

	g "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type ReceiptRepository struct {
	DB *sql.DB
}

type Receipt struct {
	UserID      string    `db:"user_id"`
	ID          string    `db:"id"`
	Total       int       `db:"total"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (r *ReceiptRepository) Create(receipt *Receipt) error {
	err := r.DB.QueryRow("INSERT INTO receipts (total, description, user_id) VALUES ($1, $2, $3) RETURNING id", receipt.Total, receipt.Description, receipt.UserID).Scan(&receipt.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (r *ReceiptRepository) CreateForUser(receipt *Receipt) error {
	err := r.DB.QueryRow("INSERT INTO receipts (total, description, user_id) VALUES ($1, $2, $3) RETURNING id", receipt.Total, receipt.Description, receipt.UserID).Scan(&receipt.ID)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func (r *ReceiptRepository) SelectAll() ([]*Receipt, error) {
	rows, err := r.DB.Query("SELECT id, total, description, created_at FROM receipts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []*Receipt
	for rows.Next() {
		receipt := &Receipt{}
		err := rows.Scan(&receipt.ID, &receipt.Total, &receipt.Description, &receipt.CreatedAt)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

func (r *ReceiptRepository) SelectForUser(userID string) ([]*Receipt, error) {
	intUserID, _ := strconv.ParseInt(userID, 10, 64)
	stmt := SELECT(
		g.Receipts.ID, g.Receipts.Total, g.Receipts.Description, g.Receipts.CreatedAt,
	).FROM(g.Receipts.Table).WHERE(g.Receipts.UserID.EQ(Int(intUserID)))

	receipts := []*Receipt{}
	err := stmt.Query(r.DB, &receipts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user receipts from DB: %w", err)
	}

	return receipts, nil
}

func (r *ReceiptRepository) FindByID(id string) (*Receipt, error) {
	receipt := &Receipt{}
	err := r.DB.QueryRow("SELECT id, total, description, created_at FROM receipts WHERE id = $1", id).Scan(&receipt.ID, &receipt.Total, &receipt.Description, &receipt.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not find receipt with id %s: %w", id, err)
	}

	return receipt, nil
}

func (r *ReceiptRepository) Delete(userID string, id string) error {
	_, err := r.DB.Exec("DELETE FROM receipts WHERE id = $1 and user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("could not delete receipt with id %s: %w", id, err)
	}

	return nil
}
