package directive

import (
	"context"
	"errors"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal/auth"
)

type GqlDirective func(ctx context.Context, obj any, next graphql.Resolver) (any, error)

func AuthDirective(tokenSecret []byte) GqlDirective {
	return func(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
		claims, ok := ctx.Value(auth.UserClaimKey).(auth.UserClaim)

		// TODO: Check the 2nd return value. Basically, handle if this failes. 1 scenario is empty token.
		if !ok || claims.State == model.UserStateGuest {
			slog.Info("no token found in context")
			return nil, errors.New("unauthorized access")
		}

		return next(ctx)
	}
}
