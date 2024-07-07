package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
	. "github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/table"
	. "github.com/go-jet/jet/v2/postgres"
)

type Receipt struct {
	model.Receipts
}

type ReceiptRepository struct {
	DB *sql.DB
}

func NewReceipt(total int32, description string, userID string) (Receipt, error) {
	receipt := Receipt{}

	receiptUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return receipt, fmt.Errorf("failed to parse user id %s : %w", userID, err)
	}

	receiptTotal := total
	receiptDescription := description
	receipt.UserID = receiptUserID
	receipt.Total = &receiptTotal
	receipt.Description = &receiptDescription

	return receipt, nil
}

func (r *ReceiptRepository) Create(receipt *Receipt) error {
	stmt := Receipts.INSERT(
		Receipts.Total, Receipts.Description, Receipts.UserID,
	).VALUES(receipt.Total, receipt.Description, receipt.UserID).RETURNING(Receipts.ID)

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return fmt.Errorf("failed to create the user: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) CreateForUser(receipt Receipt) error {
	stmt := Receipts.INSERT(
		Receipts.Total, Receipts.Description, Receipts.UserID,
	).VALUES(receipt.Total, receipt.Description, receipt.UserID).RETURNING(Receipts.ID)

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return fmt.Errorf("failed to create the user: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) SelectForUser(userID string) ([]Receipt, error) {
	stmt := SELECT(
		Receipts.ID, Receipts.Total, Receipts.Description, Receipts.CreatedAt,
	).FROM(Receipts.Table).WHERE(Receipts.UserID.EQ(RawInt(userID)))

	receipts := []Receipt{}
	err := stmt.Query(r.DB, &receipts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user receipts from DB: %w", err)
	}

	return receipts, nil
}

func (r *ReceiptRepository) FindByID(id string) (Receipt, error) {
	receipt := Receipt{}

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
