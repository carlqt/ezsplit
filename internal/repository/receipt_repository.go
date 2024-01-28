package repository

import (
	"database/sql"
	"time"
)

type ReceiptRepository struct {
	DB *sql.DB
}

type Receipt struct {
	UserID      int       `db:"user_id"`
	ID          int       `db:"id"`
	Total       int       `db:"total"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

func (r *ReceiptRepository) Create(receipt *Receipt) error {
	err := r.DB.QueryRow("INSERT INTO receipts (total, description, user_id) VALUES ($1, $2, $3) RETURNING id", receipt.Total, receipt.Description, receipt.UserID).Scan(&receipt.ID)
	if err != nil {
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
			return nil, err
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}
