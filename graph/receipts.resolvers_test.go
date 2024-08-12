package graph

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log/slog"
	"strconv"
	"testing"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
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
			claims := auth.NewUserClaim(
				user.ID,
				user.Username,
			)

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

			currentUserClaims := auth.NewUserClaim(
				receipt1.User.ID,
				receipt1.User.Username,
			)

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, currentUserClaims)

			rID := strconv.Itoa(int(receipt1.ID))

			input := &model.DeleteMyReceiptInput{
				ID: rID,
			}

			result, err := testReceiptsResolver.DeleteMyReceipt(ctx, input)

			if assert.Nil(t, err) {
				assert.Equal(t, rID, result)
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

			currentUserClaims := auth.NewUserClaim(
				receipt1.User.ID,
				receipt1.User.Username,
			)

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, currentUserClaims)

			result, err := myQueryResolver.MyReceipts(ctx)
			expectedID := strconv.Itoa(int(receipt1.ID))
			expectedUserID := strconv.Itoa(int(receipt1.User.ID))

			if assert.Nil(t, err) {
				if assert.Equal(t, 1, len(result)) {
					r := result[0]
					assert.Equal(t, "receipt 1", r.Description)
					assert.Equal(t, "99.00", r.Total)
					assert.Equal(t, expectedID, r.ID)
					assert.Equal(t, expectedUserID, r.UserID)
				}
			}
		})
	})

	t.Run("Receipt", func(t *testing.T) {
		t.Run("returns a receipt of the user", func(t *testing.T) {
			defer truncateTables()

			myQueryResolver := queryResolver{&resolvers}

			// Create user
			user, _ := app.Repositories.UserRepository.Create("sample_username", "password")

			// Create 2 receipts
			receipt1, _ := repository.NewReceipt(
				8900,
				"receipt 1",
				strconv.Itoa(int(user.ID)),
			)
			app.Repositories.ReceiptRepository.CreateForUser(&receipt1)

			// Creating a receipt for another user
			receipt2, _ := repository.NewReceipt(
				10788,
				"receipt 2",
				strconv.Itoa(int(user.ID)),
			)
			app.Repositories.ReceiptRepository.CreateForUser(&receipt2)

			id := strconv.Itoa(int(receipt2.ID))
			result, err := myQueryResolver.Receipt(context.TODO(), id)

			if assert.Nil(t, err) {
				assert.Equal(t, "receipt 2", result.Description)
				assert.Equal(t, "107.88", result.Total)
			}
		})
	})

	t.Run("GeneratePublicURL", func(t *testing.T) {
		defer truncateTables()

		// create receipt and user
		receipt, err := integration_test.CreateReceiptWithUser(app.DB, 9900, "receipt 1")
		if err != nil {
			t.Fatal(err)
		}

		currentUserClaims := auth.NewUserClaim(
			receipt.User.ID,
			receipt.User.Username,
		)

		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.UserClaimKey, currentUserClaims)

		id := strconv.Itoa(int(receipt.ID))
		result, err := testReceiptsResolver.GeneratePublicURL(ctx, id)

		hash := md5.Sum([]byte(id))
		expected := hex.EncodeToString(hash[:])

		if assert.Nil(t, err) {
			assert.Equal(t, expected, result.Slug)
			assert.Equal(t, "receipt 1", result.Description)
			assert.Equal(t, "99.00", result.Total)
		}
	})

	t.Run("RemovePublicURL", func(t *testing.T) {
		defer truncateTables()

		// create receipt and user
		receipt, err := integration_test.CreateReceiptWithUser(app.DB, 9900, "receipt 1")
		if err != nil {
			t.Fatal(err)
		}

		currentUserClaims := auth.NewUserClaim(
			receipt.User.ID,
			receipt.User.Username,
		)

		ctx := context.Background()
		ctx = context.WithValue(ctx, auth.UserClaimKey, currentUserClaims)

		id := strconv.Itoa(int(receipt.ID))
		result, err := testReceiptsResolver.RemovePublicURL(ctx, id)

		if assert.Nil(t, err) {
			assert.Equal(t, "", result.Slug)
			assert.Equal(t, "receipt 1", result.Description)
			assert.Equal(t, "99.00", result.Total)
		}
	})
}
