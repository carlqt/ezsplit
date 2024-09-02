package graph

import (
	"context"
	"testing"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	integration_test "github.com/carlqt/ezsplit/test"
	"github.com/stretchr/testify/assert"
)

func TestMeResolver(t *testing.T) {
	app := internal.NewApp()
	resolvers := Resolver{Repositories: app.Repositories, Config: app.Config}
	testMeResolver := meResolver{&resolvers}
	truncateTables := func() {
		integration_test.TruncateAllTables(app.DB)
	}

	t.Run("TotalPayables", func(t *testing.T) {
		t.Run("when there are no rows in user_orders", func(t *testing.T) {
			defer truncateTables()

			meModel := &model.Me{ID: "1"}
			result, err := testMeResolver.TotalPayables(context.TODO(), meModel)

			if assert.Nil(t, err) {
				assert.Equal(t, "0.00", result)
			}
		})
	})

	t.Run("Receipts", func(t *testing.T) {
		t.Run("when obj is nil", func(t *testing.T) {
			defer truncateTables()

			ctx := context.TODO()
			_, err := testMeResolver.Receipts(ctx, nil)

			assert.EqualError(t, err, "missing Me object")
		})
	})

	t.Run("Me", func(t *testing.T) {
		t.Run("when context is empty", func(t *testing.T) {
			queryResolver := queryResolver{&resolvers}

			_, err := queryResolver.Me(context.TODO())

			assert.Nil(t, err)
		})

		t.Run("when user exists", func(t *testing.T) {
			queryResolver := queryResolver{&resolvers}
			user, _ := app.Repositories.UserRepository.CreateGuest("john_doe")
			userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, userClaim)

			result, err := queryResolver.Me(ctx)

			if assert.Nil(t, err) {
				assert.Equal(t, "john_doe", result.Username)
				assert.Equal(t, userClaim.State, result.State)
			}
		})

		// Edge case as to when the frontend has a valid JWT and the user was deleted from the DB
		t.Run("when user does not exist", func(t *testing.T) {
			queryResolver := queryResolver{&resolvers}
			user, _ := app.Repositories.UserRepository.CreateGuest("john_doe")
			userClaim := auth.NewUserClaim(999, user.Name, user.IsVerified())

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, userClaim)

			_, err := queryResolver.Me(ctx)

			assert.ErrorContains(t, err, "can't find current user")
		})
	})
}
