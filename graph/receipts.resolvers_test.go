package graph

import (
	"context"
	"log/slog"
	"testing"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	integration_test "github.com/carlqt/ezsplit/test"
	"github.com/stretchr/testify/assert"
)

func TestReceiptsResolver(t *testing.T) {
	app := internal.NewApp()
	resolvers := Resolver{Repositories: app.Repositories, Config: app.Config}
	testReceiptsResolver := mutationResolver{&resolvers}
	truncateTables := func() {
		integration_test.TruncateAllTables(app.DB)
	}

	t.Run("CreateMyReceipt", func(t *testing.T) {
		t.Run("when userclaim exists it returns the model.Receipt struct", func(t *testing.T) {
			defer truncateTables()

			user, _ := integration_test.CreateUser(app.DB, "sample_username")
			claims := auth.UserClaim{
				ID:       user.ID,
				Username: user.Username,
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, claims)

			total := float64(255)
			input := &model.ReceiptInput{
				Description: "House Rent",
				Total:       &total,
			}

			result, err := testReceiptsResolver.CreateMyReceipt(ctx, input)

			if assert.Nil(t, err) {
				assert.Equal(t, "House Rent", result.Description)
				assert.Equal(t, claims.ID, result.UserID)
			}
		})

		t.Run("when input is empty then it returns an error", func(t *testing.T) {
			defer truncateTables()

			ctx := context.TODO()

			_, err := testReceiptsResolver.CreateMyReceipt(ctx, nil)

			assert.EqualError(t, err, "invalid input for CreateMyReceipt")
		})
	})

	t.Run("DeleteMyReceipt", func(t *testing.T) {
		t.Run("deletes a users receipt", func(t *testing.T) {
			defer truncateTables()

			receipt1, err := integration_test.CreateReceiptWithUser(app.DB, 9900, "receipt 1")
			if err != nil {
				t.Fatal(err)
			}

			integration_test.CreateReceiptWithUser(app.DB, 10000, "receipt 2")

			currentUserClaims := auth.UserClaim{
				ID:       receipt1.User.ID,
				Username: receipt1.User.Username,
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, currentUserClaims)

			input := &model.DeleteMyReceiptInput{
				ID: receipt1.ID,
			}

			result, err := testReceiptsResolver.DeleteMyReceipt(ctx, input)

			if assert.Nil(t, err) {
				assert.Equal(t, receipt1.ID, result)
			}
		})
	})

	t.Run("MyReceipts", func(t *testing.T) {
		t.Run("returns all receipts of the user", func(t *testing.T) {
			defer truncateTables()

			myQueryResolver := queryResolver{&resolvers}

			receipt1, err := integration_test.CreateReceiptWithUser(app.DB, 9900, "receipt 1")
			if err != nil {
				slog.Error(err.Error())
				t.Fatal(err)
			}

			// Creating a receipt for another user
			_, err = integration_test.CreateReceiptWithUser(app.DB, 10000, "receipt 2")
			if err != nil {
				slog.Error(err.Error())
				t.Fatal(err)
			}

			currentUserClaims := auth.UserClaim{
				ID:       receipt1.User.ID,
				Username: receipt1.User.Username,
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, currentUserClaims)

			result, err := myQueryResolver.MyReceipts(ctx)

			if assert.Nil(t, err) {
				if assert.Equal(t, 1, len(result)) {
					r := result[0]
					assert.Equal(t, "receipt 1", r.Description)
					assert.Equal(t, "99.00", r.Total)
				}
			}
		})
	})
}
