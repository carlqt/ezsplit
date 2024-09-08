package graph

import (
	"context"
	"strconv"
	"testing"

	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
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

  t.Run("Orders", func(t *testing.T) {
    t.Run("when filtering by receipt", func(t *testing.T) {
      var me model.Me
      var filter model.OrderFilterInput

      // create user
      user, err := app.Repositories.UserRepository.CreateWithAccount("jane_smith", "password")
      if err != nil {
        t.Fatal(err)
      }

      userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())
      

      ctx := context.Background()
      ctx = context.WithValue(ctx, auth.UserClaimKey, userClaim)

      // create receipt 
      receipt := repository.Receipt{}
      receipt.UserID = user.ID
      receipt.Description = "sample receipt"

      err = app.Repositories.ReceiptRepository.CreateForUser(&receipt)
      if err != nil {
        t.Fatal(err)
      }

      // create Item for receipt
      item := repository.Item{}
      item.Name = repository.Nullable("Item 1")
      item.Price = 5000
      item.ReceiptID = receipt.ID
      err = app.Repositories.ItemRepository.Create(&item)
      if err != nil {
        t.Fatal(err)
      }

      itemID := strconv.Itoa(int(item.ID))
      // assign user to item
      app.Repositories.UserOrdersRepository.Create(userClaim.ID, itemID)


      receiptID := strconv.Itoa(int(receipt.ID))
      filter.ReceiptID = receiptID
      resp, err := testMeResolver.Orders(ctx, &me, &filter)

      if assert.Nil(t, err) {
        respItem := resp[0]

        assert.Equal(t, itemID, respItem.ID)
      }
    })

    t.Run("when filter input is not nil", func(t *testing.T) {
      var me model.Me

      // create user
      user, err := app.Repositories.UserRepository.CreateWithAccount("jane_smith1", "password")
      if err != nil {
        t.Fatal(err)
      }

      userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())
      

      ctx := context.Background()
      ctx = context.WithValue(ctx, auth.UserClaimKey, userClaim)

      // create receipt 
      receipt := repository.Receipt{}
      receipt.UserID = user.ID
      receipt.Description = "sample receipt"

      err = app.Repositories.ReceiptRepository.CreateForUser(&receipt)
      if err != nil {
        t.Fatal(err)
      }

      // create Item for receipt
      item := repository.Item{}
      item.Name = repository.Nullable("Item 1")
      item.Price = 5000
      item.ReceiptID = receipt.ID
      err = app.Repositories.ItemRepository.Create(&item)
      if err != nil {
        t.Fatal(err)
      }

      itemID := strconv.Itoa(int(item.ID))

      // assign user to item
      app.Repositories.UserOrdersRepository.Create(userClaim.ID, itemID)

      resp, err := testMeResolver.Orders(ctx, &me, nil)

      if assert.Nil(t, err) {
        respItem := resp[len(resp) - 1]

        assert.Equal(t, itemID, respItem.ID)
      }
    })
  })
}
