package graph

import (
	"math/big"

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
	return &model.Receipt{
		ID:          string(receipt.ID),
		Total:       toPriceDisplay(*receipt.Total),
		Description: *receipt.Description,
		UserID:      string(receipt.UserID),
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
