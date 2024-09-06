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
      _, err := testMutationResolver.AssignOrRemoveMeFromItem(ctx, itemID)

      assert.Nil(t, err)
    })

    t.Run("when user is assigned to an item", func(t *testing.T) {
    })
  })
}
