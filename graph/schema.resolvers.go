package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"context"
	"fmt"
	"strconv"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/repository"
)

// CreateReceipt is the resolver for the createReceipt field.
func (r *mutationResolver) CreateReceipt(ctx context.Context, input *model.ReceiptInput) (*model.Receipt, error) {
	receipt := &repository.Receipt{
		Total:       int(*input.Price * 100),
		Description: input.Description,
		// TODO: UserID should come from the context.
		// UserID: 1,
	}

	err := r.Repositories.ReceiptRepository.Create(receipt)
	if err != nil {
		return nil, err
	}

	receiptResponse := &model.Receipt{
		ID:          strconv.Itoa(receipt.ID),
		Total:       input.Price,
		Description: input.Description,
	}

	return receiptResponse, nil
}

// AddItemToReceipt is the resolver for the addItemToReceipt field.
func (r *mutationResolver) AddItemToReceipt(ctx context.Context, input *model.AddItemToReceiptInput) (*model.Item, error) {
	receiptID, err := strconv.Atoi(input.ReceiptID)
	if err != nil {
		return nil, err
	}

	// Initialize repository struct using the input
	item := repository.Item{
		ReceiptID: receiptID,
		Name:      input.Name,
		Price:     int(*input.Price * 100),
	}

	err = r.Repositories.ItemRepository.Create(&item)
	if err != nil {
		return nil, err
	}

	return &model.Item{
		ID:    strconv.Itoa(item.ID),
		Name:  item.Name,
		Price: input.Price,
	}, nil
}

// AssignUserToItem is the resolver for the assignUserToItem field.
func (r *mutationResolver) AssignUserToItem(ctx context.Context, input *model.AssignUserToItemInput) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: AssignUserToItem - assignUserToItem"))
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.UserInput) (*model.User, error) {
	user := repository.User{
		Username: input.Username,
	}

	err := r.Repositories.UserRepository.Create(&user)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:       strconv.Itoa(user.ID),
		Username: user.Username,
	}, nil
}

// GetReceipts is the resolver for the getReceipts field.
func (r *queryResolver) GetReceipts(ctx context.Context) ([]*model.Receipt, error) {
	receipts, err := r.Repositories.ReceiptRepository.SelectAll()
	if err != nil {
		return nil, err
	}

	var modelReceipts []*model.Receipt

	for _, receipt := range receipts {
		total := float64(receipt.Total) / 100

		modelReceipt := &model.Receipt{
			ID:          strconv.Itoa(receipt.ID),
			Total:       &total,
			Description: receipt.Description,
		}
		modelReceipts = append(modelReceipts, modelReceipt)
	}

	return modelReceipts, nil
}

// GetReceiptByID is the resolver for the getReceiptById field.
func (r *queryResolver) GetReceiptByID(ctx context.Context) (*model.Receipt, error) {
	panic(fmt.Errorf("not implemented: GetReceiptByID - getReceiptById"))
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.Repositories.UserRepository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var modelUsers []*model.User
	for _, user := range users {
		modelUser := &model.User{
			ID:       strconv.Itoa(user.ID),
			Username: user.Username,
		}
		modelUsers = append(modelUsers, modelUser)
	}

	return modelUsers, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
