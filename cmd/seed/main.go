package main

import (
	"log/slog"

	"github.com/carlqt/ezsplit/.gen/ezsplit_dev/public/model"
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

	userClaim := auth.NewUserClaim(user.ID, user.Username)
	signedToken, err := auth.CreateAndSignToken(userClaim, secret)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("User created", "accessToken", signedToken)
	return userClaim.ID, nil
}

func createReceipt(repo *repository.ReceiptRepository, userID string) (string, error) {
	receipt, _ := repository.NewReceipt(
		4000,
		"Jollibee",
		userID,
	)

	err := repo.Create(&receipt)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("Receipt created", "receiptID", receipt.ID)
	return string(receipt.ID), nil
}

func createItems(repo *repository.ItemRepository, receiptID string) error {
	price := int32(4000)
	name := "Chickenjoy"
	items := []repository.Item{
		repository.Item{
			model.Items{
				Name: &name, Price: &price, ReceiptID: repository.BigInt(receiptID),
			},
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
