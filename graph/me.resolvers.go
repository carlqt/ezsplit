package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/auth"
)

// TotaylPayables is the resolver for the totaylPayables field.
func (r *meResolver) TotaylPayables(ctx context.Context, obj *model.Me) (string, error) {
	// panic(fmt.Errorf("not implemented: TotaylPayables - totaylPayables"))
	return "222", nil
}

// Me is the resolver for the Me field.
func (r *queryResolver) Me(ctx context.Context) (*model.Me, error) {
	claims := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

	user, err := r.Repositories.UserRepository.FindByID(claims.ID)
	if err != nil {
		return nil, err
	}

	return &model.Me{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}

// Me returns MeResolver implementation.
func (r *Resolver) Me() MeResolver { return &meResolver{r} }

type meResolver struct{ *Resolver }
