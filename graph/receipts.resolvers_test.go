package graph

import (
	"context"
	"database/sql"
	"log/slog"
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
		defer truncateTables()
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
		defer truncateTables()
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

			// CreateReceiptWithRandomUser(app.DB, 8888, "McDonalds")
			// CreateReceiptWithRandomUser(app.DB, 8888, "McDonalds")
			// CreateReceiptWithRandomUser(app.DB, 8888, "McDonalds")
			// CreateReceiptWithRandomUser(app.DB, 8888, "McDonalds")
			// CreateReceiptWithRandomUser(app.DB, 8888, "McDonalds")

			// integration_test.CreateReceiptWithUser(app.DB, 9999, "McDondalds")
			// integration_test.CreateReceiptWithUser(app.DB, 9999, "McDondalds")
			// integration_test.CreateReceiptWithUser(app.DB, 9999, "McDondalds")
			// integration_test.CreateReceiptWithUser(app.DB, 9999, "McDondalds")
			// integration_test.CreateReceiptWithUser(app.DB, 9999, "McDondalds")
			// integration_test.CreateReceiptWithUser(app.DB, 9999, "McDondalds")

			// usersCount(app.DB)

			myQueryResolver := queryResolver{&resolvers}

			receipt1, err := integration_test.CreateReceiptWithUser(app.DB, 9900, "receipt 1")
			if err != nil {
				slog.Error(err.Error())
				t.Fatal(err)
			}

			// usersCount(app.DB)
			// checkReceipts(app.DB, receipt1.UserID)

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

			usersCount(app.DB)
			checkReceipts(app.DB, receipt1.UserID)

			if assert.Nil(t, err) {
				if assert.Equal(t, 1, len(result), result) {
					r := result[0]
					assert.Equal(t, "receipt 1", r.Description)
					assert.Equal(t, "99.00", r.Total)
				}
			}
		})
	})
}

func checkReceipts(db *sql.DB, userID string) {
	var receipts []*repository.Receipt

	// slog.Debug("fetching userID", "userID", userID)

	rows, err := db.Query("select id, user_id, description, total from receipts WHERE user_id = $1", userID)
	if err != nil {
		slog.Error("failed to fetch all receipts", "err", err)
	}

	for rows.Next() {
		r := repository.Receipt{}

		rows.Scan(&r.ID, &r.UserID, &r.Description, &r.Total)
		receipts = append(receipts, &r)
	}

	slog.Debug("receipt checks", "receipts", receipts)
}

func checkUsers(db *sql.DB) {
	rows, err := db.Query("select id, username from users")
	if err != nil {
		slog.Error("failed to fetch all users", "err", err)
	}

	users := make([]repository.User, 0)

	for rows.Next() {
		u := repository.User{}

		rows.Scan(&u.ID, &u.Username)
		users = append(users, u)
	}

	slog.Debug("users check", "users", users)
}

func usersCount(db *sql.DB) {
	count := 0

	err := db.QueryRow("select count(*) from users").Scan(&count)
	if err != nil {
		slog.Error("failed to count users", "err", err)
	}

	slog.Debug("User count", "count", count)
}
