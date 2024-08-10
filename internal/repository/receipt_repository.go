package repository

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/carlqt/ezsplit/.gen/public/model"
	. "github.com/carlqt/ezsplit/.gen/public/table"
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
	receipt.UserID = receiptUserID
	receipt.Total = &receiptTotal
	receipt.Description = description

	return receipt, nil
}

func (r *ReceiptRepository) Create(receipt *Receipt) error {
	stmt := Receipts.INSERT(
		Receipts.Total, Receipts.Description, Receipts.UserID,
	).VALUES(receipt.Total, receipt.Description, receipt.UserID).RETURNING(Receipts.AllColumns)

	err := stmt.Query(r.DB, receipt)
	if err != nil {
		return fmt.Errorf("failed to create the user: %w", err)
	}

	return nil
}

func (r *ReceiptRepository) CreateForUser(receipt *Receipt) error {
	stmt := Receipts.INSERT(
		Receipts.Total, Receipts.Description, Receipts.UserID,
	).VALUES(receipt.Total, receipt.Description, receipt.UserID).RETURNING(Receipts.AllColumns)

	err := stmt.Query(r.DB, receipt)
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
	).FROM(Receipts.Table).WHERE(Receipts.ID.EQ(RawInt(id)))

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return receipt, fmt.Errorf("could not find receipt with id %s: %w", id, err)
	}

	return receipt, nil
}

func (r *ReceiptRepository) Delete(userID string, id string) error {
	stmt := Receipts.DELETE().WHERE(Receipts.ID.EQ(RawInt(id)).AND(Receipts.UserID.EQ(RawInt(userID))))

	_, err := stmt.Exec(r.DB)
	if err != nil {
		return fmt.Errorf("could not delete receipt with id %s: %w", id, err)
	}

	return nil
}

func (r *ReceiptRepository) GeneratePublicUrlPath(userID, receiptID string) (Receipt, error) {
	var receipt Receipt

	hash := md5.Sum([]byte(receiptID))
	publicUrlPath := hex.EncodeToString(hash[:])
	receipt.PublicURLPath = publicUrlPath

	stmt := Receipts.UPDATE(
		Receipts.PublicURLPath,
	).MODEL(receipt).WHERE(
		Receipts.UserID.EQ(RawInt(userID)).AND(Receipts.ID.EQ(RawInt(receiptID))),
	).RETURNING(Receipts.AllColumns)

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return receipt, fmt.Errorf("failed to update receipt with id=%s: %w", receiptID, err)
	}

	return receipt, nil
}

// Update method is mainly used to easily create data for tests
func (r *ReceiptRepository) Update(receipt *Receipt) error {
	stmt := Receipts.UPDATE(
		Receipts.MutableColumns,
	).MODEL(receipt).WHERE(
		Receipts.ID.EQ(Int32(receipt.ID)),
	).RETURNING(Receipts.AllColumns)

	err := stmt.Query(r.DB, &receipt)
	if err != nil {
		return fmt.Errorf("failed to update receipt %v: %w", receipt, err)
	}

	return nil
}
