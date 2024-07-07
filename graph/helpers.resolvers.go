package graph

import (
	"math/big"
	"strconv"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/repository"
)

func newModelUser(user *repository.User) *model.User {
	return &model.User{
		ID:       user.ID,
		Username: user.Username,
	}
}

// newModelItem is a constructor that converts a repository.Item to a model.Item
func newModelItem(item *repository.Item) *model.Item {
	price := toPriceDisplay(item.Price)

	return &model.Item{
		ID:    item.ID,
		Name:  item.Name,
		Price: price,
	}
}

func newModelReceipt(receipt *repository.Receipt) *model.Receipt {
	userID := strconv.Itoa(int(receipt.UserID))
	receiptID := strconv.Itoa(int(receipt.ID))

	return &model.Receipt{
		ID:          receiptID,
		Total:       toPriceDisplay(*receipt.Total),
		Description: *receipt.Description,
		UserID:      userID,
	}
}

// toPriceCents converts an input price to cents (lowest denomination).
// An explanation can be found on https://stackoverflow.com/a/46492064/976035
func toPriceCents(price float64) int32 {
	return int32(price*100 + 0.5)
}

// toPriceDisplay converts prices to a displayable format (e.g. 1000 -> 10.00)
func toPriceDisplay[T int | int32](price T) string {
	p := big.NewRat(int64(price), 100)

	return p.FloatString(2)
}
