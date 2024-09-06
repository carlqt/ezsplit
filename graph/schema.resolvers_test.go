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

func TestSchemaResolver(t *testing.T) {
	app := internal.NewApp()
	resolvers := Resolver{Repositories: app.Repositories, Config: app.Config}
	testItemResolver := itemResolver{&resolvers}
  testMutationResolver := mutationResolver{&resolvers}
  
	truncateTables := func() {
		integration_test.TruncateAllTables(app.DB)
	}

	t.Run("SharedBy", func(t *testing.T) {
		t.Run("when there are no results", func(t *testing.T) {
			defer truncateTables()

			modelItem := model.Item{ID: "888"}
			result, _ := testItemResolver.SharedBy(context.TODO(), &modelItem)

			assert.Empty(t, result, result)
			assert.Zero(t, len(result))
		})
	})

  t.Run("AssignOrRemoveMeFromItem", func(t *testing.T) {
    t.Run("when user is not assigned to an item", func(t *testing.T) {
      // create user
      user, _ := app.Repositories.UserRepository.CreateWithAccount("john_doe", "password")
			userClaim := auth.NewUserClaim(user.ID, user.Name, user.IsVerified())

			ctx := context.Background()
			ctx = context.WithValue(ctx, auth.UserClaimKey, userClaim)

      // create receipt 
      receipt := repository.Receipt{}
      receipt.UserID = repository.BigInt(user.ID)
      receipt.Description = "sample receipt"

      app.Repositories.ReceiptRepository.CreateForUser(&receipt)

      // create Item for receipt
      item := repository.Item{}
      item.Name = repository.Nullable("Item 1")
      item.Price = 5000
      item.ReceiptID = repository.BigInt(receipt.ID)
      app.Repositories.ItemRepository.Create(&item)

      itemID := strconv.Itoa(int(item.ID))
      resp, err := testMutationResolver.AssignOrRemoveMeFromItem(ctx, itemID)

      if assert.Nil(t, err) {
        assert.Equal(t, itemID, resp.ItemID)
        assert.Equal(t, userClaim.ID, resp.UserID)
      }
    })

    t.Run("when user is assigned to an item", func(t *testing.T) {
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
      receipt.UserID = repository.BigInt(user.ID)
      receipt.Description = "sample receipt"

      err = app.Repositories.ReceiptRepository.CreateForUser(&receipt)
      if err != nil {
        t.Fatal(err)
      }

      // create Item for receipt
      item := repository.Item{}
      item.Name = repository.Nullable("Item 1")
      item.Price = 5000
      item.ReceiptID = repository.BigInt(receipt.ID)
      err = app.Repositories.ItemRepository.Create(&item)
      if err != nil {
        t.Fatal(err)
      }

      itemID := strconv.Itoa(int(item.ID))

      // assign user to item
      app.Repositories.UserOrdersRepository.Create(userClaim.ID, itemID)

      resp, err := testMutationResolver.AssignOrRemoveMeFromItem(ctx, itemID)

      if assert.Nil(t, err) {
        assert.Equal(t, itemID, resp.ItemID)
        assert.Equal(t, userClaim.ID, resp.UserID)
      }
    })
  })
}
