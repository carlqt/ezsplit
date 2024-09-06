package graph

import (
	"math/big"
	"strconv"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/repository"
)

func newUserOrderRef(userID, itemID string) *model.UserOrderRef {
  return &model.UserOrderRef{
    ItemID: itemID,
    UserID: userID,
  }
}

func newModelUser(userID int32, username string, isVerified bool) *model.User {
	id := strconv.Itoa(int(userID))
	state := model.UserStateGuest

	if isVerified {
		state = model.UserStateVerified
	}

	return &model.User{
		ID:       id,
		Username: username,
		State:    state,
	}
}

func newModelMe(userID int32, username string, isVerified bool) *model.Me {
	id := strconv.Itoa(int(userID))
	state := model.UserStateGuest

	if isVerified {
		state = model.UserStateVerified
	}

	return &model.Me{
		ID:       id,
		Username: username,
		State:    state,
	}
}

func newModelUserWithJwt(userID int32, username string, accessToken string) *model.UserWithJwt {
	id := strconv.Itoa(int(userID))

	return &model.UserWithJwt{
		ID:          id,
		Username:    username,
		AccessToken: accessToken,
	}
}

// newModelItem is a constructor that converts a repository.Item to a model.Item
func newModelItem(item repository.Item) *model.Item {
	price := toPriceDisplay(item.Price)
	id := strconv.Itoa(int(item.ID))

	return &model.Item{
		ID:    id,
		Name:  *item.Name,
		Price: price,
	}
}

func newModelReceipt(receipt *repository.Receipt) *model.Receipt {
	userID := strconv.Itoa(int(receipt.UserID))
	receiptID := strconv.Itoa(int(receipt.ID))

	return &model.Receipt{
		ID:          receiptID,
		Total:       toPriceDisplay(*receipt.Total),
		Description: receipt.Description,
		Slug:        receipt.URLSlug,
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
