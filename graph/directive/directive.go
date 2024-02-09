package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

type GqlDirective func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error)

func AuthDirective(tokenSecret []byte) GqlDirective {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		// bearerToken := ctx.Value(auth.TokenKey).(string)

		// claims, err := auth.ValidateBearerToken(bearerToken, tokenSecret)
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	return nil, errors.New("unauthorized access")
		// }

		// ctx = context.WithValue(ctx, auth.UserClaimKey, claims)

		return next(ctx)
	}
}
