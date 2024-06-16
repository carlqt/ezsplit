package main

import (
	"log/slog"

	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
)

func main() {
	app := internal.NewApp()

	userID, err := createUser(app.Repositories.UserRepository, app.Config.JWTSecret)
	if err != nil {
		return
	}

	receiptID, err := createReceipt(app.Repositories.ReceiptRepository, userID)
	if err != nil {
		return
	}

	_ = createItems(app.Repositories.ItemRepository, receiptID)
}

func createUser(repo *repository.UserRepository, secret []byte) (string, error) {
	password := "password"
	hashedPassword, _ := auth.HashPassword(password)

	user, err := repo.Create("john_smith", hashedPassword)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	userClaim := auth.UserClaim{
		ID:       user.ID,
		Username: user.Username,
	}
	signedToken, err := auth.CreateAndSignToken(userClaim, secret)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("User created", "accessToken", signedToken)
	return user.ID, nil
}

func createReceipt(repo *repository.ReceiptRepository, userID string) (string, error) {
	receipt := repository.Receipt{
		UserID:      userID,
		Description: "Jollibee",
		Total:       4000,
	}

	err := repo.Create(&receipt)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("Receipt created", "receiptID", receipt.ID)
	return receipt.ID, nil
}

func createItems(repo *repository.ItemRepository, receiptID string) error {
	items := []repository.Item{
		{
			ReceiptID: receiptID,
			Name:      "Chickenjoy",
			Price:     4000,
		},
	}

	for _, item := range items {
		err := repo.Create(&item)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	slog.Info("Items created")
	return nil
}
