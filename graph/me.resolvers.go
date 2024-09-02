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
)

// TotalPayables is the resolver for the totalPayables field.
func (r *meResolver) TotalPayables(ctx context.Context, obj *model.Me) (string, error) {
	// panic(fmt.Errorf("not implemented: TotalPayables - totalPayables"))
	totalPayables, err := r.Repositories.UserOrdersRepository.GetTotalPayables(obj.ID)
	if err != nil {
		slog.Error(err.Error())
		return "", errors.New("failed to get the total payables")
	}

	return toPriceDisplay(totalPayables), nil
}

// Receipts is the resolver for the receipts field.
func (r *meResolver) Receipts(ctx context.Context, obj *model.Me) ([]*model.Receipt, error) {
	if obj == nil {
		return nil, errors.New("missing Me object")
	}

	receipts, err := r.Repositories.ReceiptRepository.SelectForUser(obj.ID)
	if err != nil {
		slog.Debug(err.Error())
		return nil, err
	}

	modelReceipts := make([]*model.Receipt, 0)

	for _, receipt := range receipts {
		modelReceipt := newModelReceipt(&receipt)
		modelReceipts = append(modelReceipts, modelReceipt)
	}

	return modelReceipts, nil
}

// Me is the resolver for the Me field.
func (r *queryResolver) Me(ctx context.Context) (*model.Me, error) {
	// Reminder for me. We can trust the claims in the context.
	// If the claims have been tampered in the frontend, it will be caught by the JWTMiddleware
	claims, ok := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	// !ok means the request didn't have a cookie or
	// The cookie didn't have a bearerToken field
	if !ok {
		return nil, nil
	}

	user, err := r.Repositories.UserRepository.FindByID(claims.ID)
	if err != nil {
		slog.Debug(err.Error())
		return nil, errors.New("can't find current user")
	}

	return newModelMe(user.ID, user.Name, user.IsVerified()), nil
}

// Me returns MeResolver implementation.
func (r *Resolver) Me() MeResolver { return &meResolver{r} }

type meResolver struct{ *Resolver }
