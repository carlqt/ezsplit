package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
	. "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type ReceiptRepository struct {
	DB *sql.DB
}

func (r *ReceiptRepository) Create(receipt model.Receipts) error {
	stmt := Receipts.INSERT(
		Receipts.Total, Receipts.Description, Receipts.UserID,
	).VALUES(receipt.Total, receipt.Description, receipt.UserID).RETURNING(Receipts.ID)

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return fmt.Errorf("failed to create the user: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) CreateForUser(receipt model.Receipts) error {
	stmt := Receipts.INSERT(
		Receipts.Total, Receipts.Description, Receipts.UserID,
	).VALUES(receipt.Total, receipt.Description, receipt.UserID).RETURNING(Receipts.ID)

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return fmt.Errorf("failed to create the user: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) SelectForUser(userID string) ([]model.Receipts, error) {
	stmt := SELECT(
		Receipts.ID, Receipts.Total, Receipts.Description, Receipts.CreatedAt,
	).FROM(Receipts.Table).WHERE(Receipts.UserID.EQ(RawInt(userID)))

	receipts := []model.Receipts{}
	err := stmt.Query(r.DB, &receipts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user receipts from DB: %w", err)
	}

	return receipts, nil
}

func (r *ReceiptRepository) FindByID(id string) (model.Receipts, error) {
	receipt := model.Receipts{}

	stmt := SELECT(
		Receipts.ID, Receipts.Total, Receipts.Description, Receipts.CreatedAt,
	).FROM(Receipts.Table).WHERE(Receipts.UserID.EQ(RawInt(id)))

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return receipt, fmt.Errorf("could not find receipt with id %s: %w", id, err)
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
