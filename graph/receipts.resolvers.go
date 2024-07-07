package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"log/slog"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
)

// CreateMyReceipt is the resolver for the createMyReceipt field.
func (r *mutationResolver) CreateMyReceipt(ctx context.Context, input *model.ReceiptInput) (*model.Receipt, error) {
	if input == nil {
		return nil, errors.New("invalid input for CreateMyReceipt")
	}

	userClaim := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	receipt, err := repository.NewReceipt(
		toPriceCents(*input.Total),
		input.Description,
		userClaim.ID,
	)

	if err != nil {
		slog.Error(err.Error(), "userID", userClaim.ID)
		return nil, errors.New("failed to read input")
	}

	err = r.Repositories.ReceiptRepository.CreateForUser(receipt)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.New("failed to create user")
	}

	return newModelReceipt(&receipt), nil
}

// DeleteMyReceipt is the resolver for the deleteMyReceipt field.
func (r *mutationResolver) DeleteMyReceipt(ctx context.Context, input *model.DeleteMyReceiptInput) (string, error) {
	userClaim := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	err := r.Repositories.ReceiptRepository.Delete(userClaim.ID, input.ID)
	if err != nil {
		slog.Error(err.Error())
		return "", errors.New("failed to delete receipt")
	}

	return input.ID, nil
}

// MyReceipts is the resolver for the myReceipts field.
func (r *queryResolver) MyReceipts(ctx context.Context) ([]*model.Receipt, error) {
	userClaim := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	receipts, err := r.Repositories.ReceiptRepository.SelectForUser(userClaim.ID)

	if err != nil {
		slog.Debug(err.Error())
		return nil, errors.New("couldn't fetch receipts of user")
	}

	modelReceipts := make([]*model.Receipt, 0)

	for _, receipt := range receipts {
		modelReceipt := newModelReceipt(receipt)
		modelReceipts = append(modelReceipts, modelReceipt)
	}

	return modelReceipts, nil
}

// Receipt is the resolver for the receipt field.
func (r *queryResolver) Receipt(ctx context.Context, id string) (*model.Receipt, error) {
	receipt, err := r.Repositories.ReceiptRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return newModelReceipt(receipt), nil
}

// User is the resolver for the user field.
func (r *receiptResolver) User(ctx context.Context, obj *model.Receipt) (*model.User, error) {
	user, err := r.Repositories.UserRepository.FindByID(obj.UserID)
	if err != nil {
		// log.Println(err)
		return nil, err
	}

	return &model.User{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// Items is the resolver for the items field.
func (r *receiptResolver) Items(ctx context.Context, obj *model.Receipt) ([]*model.Item, error) {
	items, err := r.Repositories.ItemRepository.SelectAllForReceipt(obj.ID)
	if err != nil {
		return nil, err
	}

	var modelItems []*model.Item
	for _, item := range items {
		modelItem := newModelItem(item)
		modelItems = append(modelItems, modelItem)
	}

	return modelItems, nil
}

// Receipt returns ReceiptResolver implementation.
func (r *Resolver) Receipt() ReceiptResolver { return &receiptResolver{r} }

type receiptResolver struct{ *Resolver }
