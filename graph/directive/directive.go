package directive

import (
	"context"
	"errors"
	"log/slog"

	"github.com/99designs/gqlgen/graphql"
	"github.com/carlqt/ezsplit/internal/auth"
)

type GqlDirective func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)

func AuthDirective(tokenSecret []byte) GqlDirective {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		// TODO: Check the 2nd return value. Basically, handle if this failes. 1 scenario is empty token.
		bearerToken, ok := ctx.Value(auth.TokenKey).(string)
		if !ok {
			slog.Error("no token found in context")
			return nil, errors.New("unauthorized access")
		}

		claims, err := auth.ValidateBearerToken(bearerToken, tokenSecret)
		if err != nil {
			slog.Error(err.Error())
			return nil, errors.New("unauthorized access")
		}

		ctx = context.WithValue(ctx, auth.UserClaimKey, claims)

		return next(ctx)
	}
}
