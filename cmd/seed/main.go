package main

import (
	"log/slog"
	"strconv"

	"github.com/carlqt/ezsplit/.gen/public/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
)

func main() {
	app := internal.NewApp()

	userID, err := createUser(app.Repositories.UserRepository, app.Config.JWTSecret)
	if err != nil {
		panic(err)
	}

	receiptID, err := createReceipt(app.Repositories.ReceiptRepository, userID)
	if err != nil {
		panic(err)
	}

	_ = createItems(app.Repositories.ItemRepository, receiptID)
}

func createUser(repo *repository.UserRepository, secret []byte) (string, error) {
	user, err := repo.CreateWithAccount("john_smith", "password")
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())
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

	err := repo.CreateForUser(&receipt)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	slog.Info("Receipt created", "receiptID", receipt.ID)

	receiptID := strconv.Itoa(int(receipt.ID))
	return receiptID, nil
}

func createItems(repo *repository.ItemRepository, receiptID string) error {
	itemsData := []struct {
		price int32
		name  string
	}{
		{4000, "Chickenjoy"},
		{2000, "Spaghetti"},
		{1000, "Burger Steak"},
	}

	for _, i := range itemsData {
		item := repository.Item{
			Items: model.Items{
				Name: &i.name, Price: i.price, ReceiptID: repository.BigInt(receiptID),
			},
		}

		err := repo.Create(&item)
		if err != nil {
			slog.Error(err.Error())
			return err
		}
	}

	slog.Info("Items created")
	return nil
}
