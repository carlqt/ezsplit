package integration_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	g "github.com/carlqt/ezsplit/.gen/public/model"
	"github.com/carlqt/ezsplit/graph"
	"github.com/carlqt/ezsplit/graph/directive"
	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestResolvers(t *testing.T) {
	// Setting up the server
	app := internal.NewApp()
	resolvers := &graph.Resolver{Repositories: app.Repositories, Config: app.Config}
	config := graph.Config{Resolvers: resolvers}
	config.Directives.Authenticated = directive.AuthDirective(app.Config.JWTSecret)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(config))
	c := client.New(internal.BearerTokenMiddleware(internal.InjectSetCookieMiddleware(srv)))

	toString := func(i int32) string {
		return strconv.Itoa(int(i))
	}

	defer TruncateAllTables(app.DB)

	t.Run("loginUser mutation", func(t *testing.T) {
		t.Run("when there are no inputs", func(t *testing.T) {
			defer TruncateAllTables(app.DB)

			_, err := CreateVerifiedUser(app.DB, "mutation_user160")
			if err != nil {
				t.Fatal(err)
			}

			query := `mutation loginUser {
				loginUser {
					username
					id
				}
			}`

			var resp struct {
				LoginUser struct {
					Username string
					Id       string
				}
			}

			err = c.Post(query, &resp)

			if assert.NotNil(t, err) {
				assert.EqualError(t, err, `[{"message":"incorrect username or password","path":["loginUser"]}]`)
			}
		})

		t.Run("when password is correct", func(t *testing.T) {
			defer TruncateAllTables(app.DB)

			_, err := CreateVerifiedUser(app.DB, "mutation_user160")
			if err != nil {
				t.Fatal(err)
			}

			query := `mutation loginUser {
				loginUser(input: {username: "mutation_user160", password: "password"}) {
					username
					id
				}
			}`

			var resp struct {
				LoginUser struct {
					Username string
					Id       string
				}
			}

			err = c.Post(query, &resp)

			if assert.Nil(t, err, resp.LoginUser) {
				assert.Equal(t, "mutation_user160", resp.LoginUser.Username)
			}
		})

		t.Run("when password is wrong", func(t *testing.T) {
			defer TruncateAllTables(app.DB)

			_, err := CreateVerifiedUser(app.DB, "mutation_user160")
			if err != nil {
				t.Fatal(err)
			}

			query := `mutation loginUser {
				loginUser(input: {username: "mutation_user160", password: "passwordX"}) {
					username
					id
				}
			}`

			var resp struct {
				LoginUser struct {
					Username string
					Id       string
				}
			}

			err = c.Post(query, &resp)

			if assert.NotNil(t, err) {
				assert.EqualError(t, err, `[{"message":"incorrect username or password","path":["loginUser"]}]`)
			}
		})
	})

	t.Run("createUser mutation", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		query := `mutation createUser {
			createUser(input: {username: "mutation_user160", password: "password", confirmPassword: "password" }) {
				username
				id
				accessToken
			}
		}`

		var resp struct {
			CreateUser struct {
				Username    string
				Id          string
				AccessToken string
			}
		}

		err := c.Post(query, &resp)

		if assert.Nil(t, err) {
			assert.Equal(t, "mutation_user160", resp.CreateUser.Username)
		}
	})

	t.Run("logoutUser mutation works", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		query := `mutation logoutUser {
			logoutUser
		}`

		var resp struct {
			LogoutUser string
		}

		err := c.Post(query, &resp)

		if assert.Nil(t, err) {
			assert.Equal(t, "ok", resp.LogoutUser)
		}
	})

	t.Run("query Me", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		user, err := CreateVerifiedUser(app.DB, "fake_user")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("with Receipts field", func(t *testing.T) {
			userClaim := auth.NewUserClaim(
				user.ID,
				user.Username,
				user.State,
			)
			accessToken, err := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)
			if err != nil {
				t.Fatal(err)
			}

			t.Run("when a receipt exists", func(t *testing.T) {
				t.Cleanup(func() {
					_, err := app.DB.Exec("DELETE FROM receipts")
					if err != nil {
						t.Fatal(err)
					}
				})

				userID := strconv.Itoa(int(user.ID))
				receipt, _ := repository.NewReceipt(35000, "test receipt", userID)
				err = app.Repositories.ReceiptRepository.CreateForUser(&receipt)
				if err != nil {
					t.Fatal(err)
				}

				query := `query Me {
          me {
            username
            id
            receipts {
              id
              description
              total
            }
          }
        }`

				var resp struct {
					Me struct {
						Username string
						Id       string
						Receipts []*model.Receipt
					}
				}

				option := func(bd *client.Request) {
					bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
				}

				err = c.Post(query, &resp, option)

				if assert.Nil(t, err) {
					responseReceipt := resp.Me.Receipts[0]
					expectedReceiptID := strconv.Itoa(int(receipt.ID))

					assert.Equal(t, user.Username, resp.Me.Username)
					assert.Equal(t, expectedReceiptID, responseReceipt.ID)
					assert.Equal(t, "test receipt", responseReceipt.Description)
					assert.Equal(t, "350.00", responseReceipt.Total)

					// Formatted total
					assert.Equal(t, "350.00", responseReceipt.Total)
				}
			})

			t.Run("when receipts is empty", func(t *testing.T) {
				query := `query Me {
          me {
            username
            id
            receipts {
              id
              description
              total
            }
          }
        }`

				var resp struct {
					Me struct {
						Username string
						Id       string
						Receipts []*model.Receipt
					}
				}

				option := func(bd *client.Request) {
					bd.HTTP.AddCookie(&http.Cookie{Name: internal.BearerTokenCookie, Value: accessToken})
				}

				err = c.Post(query, &resp, option)

				if assert.Nil(t, err) {
					responseReceipt := resp.Me.Receipts

					assert.Equal(t, userClaim.ID, resp.Me.Id)
					assert.Empty(t, responseReceipt)
				}
			})
		})

		t.Run("when jwt exists", func(t *testing.T) {
			userClaim := auth.NewUserClaim(
				user.ID,
				user.Username,
				user.State,
			)
			accessToken, err := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)
			if err != nil {
				t.Fatal(err)
			}

			query := `query Me {
				me {
					username
					id
				}
			}`

			var resp struct {
				Me struct {
					Username string
					Id       string
				}
			}

			option := func(bd *client.Request) {
				bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
			}

			err = c.Post(query, &resp, option)

			if assert.Nil(t, err) {
				assert.Equal(t, user.Username, resp.Me.Username)
				assert.Equal(t, userClaim.ID, resp.Me.Id)
			}
		})

		t.Run("when jwt does not exists", func(t *testing.T) {
			query := `query Me {
				me {
					username
					id
				}
			}`

			var resp struct {
				Me struct {
					Username string
					Id       string
				}
			}

			err = c.Post(query, &resp)

			if assert.NotNil(t, err) {
				// TODO: There should be a better way to check the error message
				assert.EqualError(t, err, `[{"message":"unauthorized access","path":["me"]}]`)
			}
		})
	})

	t.Run("mutation createMyReceipt", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		user, err := CreateVerifiedUser(app.DB, "fake_user")
		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := user.GetAuthToken(app.Config.JWTSecret)
		if err != nil {
			t.Fatal(err)
		}

		query := `mutation createMyReceipt {
			createMyReceipt(input: {description: "test receipt", total: 350 }) {
				description
				total
				id
			}
		}`

		var resp struct {
			CreateMyReceipt struct {
				Description string
				Total       string
				Id          string
			}
		}

		option := func(bd *client.Request) {
			bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.Equal(t, "test receipt", resp.CreateMyReceipt.Description)
			assert.Equal(t, "350.00", resp.CreateMyReceipt.Total)
		}
	})

	t.Run("mutation assignMeToItem", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		user, err := app.Repositories.UserRepository.Create("john_doe", "testing")
		if err != nil {
			t.Fatal(err)
		}

		userID := strconv.Itoa(int(user.ID))
		receipt, _ := repository.NewReceipt(35000, "test receipt", userID)
		err = app.Repositories.ReceiptRepository.UnsafeCreate(&receipt)
		if err != nil {
			t.Fatal(err)
		}

		itemName := "Dumplings"
		itemPrice := int32(10000)

		item := repository.Item{} //{model.Items{Name: "Dumplings", Price: 10000, ReceiptID: itemReceiptID}
		item.Name = &itemName
		item.Price = itemPrice
		item.ReceiptID = repository.BigInt(receipt.ID)

		err = app.Repositories.ItemRepository.Create(&item)
		if err != nil {
			t.Fatal(err)
		}

		userClaim := auth.NewUserClaim(
			user.ID,
			user.Username,
			user.State,
		)
		accessToken, _ := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

		query := fmt.Sprintf(`mutation assignMeToItem{
			assignMeToItem(input: { itemId: "%d" }) {
				name
				price
				id
				sharedBy {
					id
					username
				}
			}
		}`, item.ID)

		var resp struct {
			AssignMeToItem struct {
				Price    string
				Name     string
				Id       string
				SharedBy []struct {
					Id       string
					Username string
				}
			}
		}

		option := func(bd *client.Request) {
			bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			itemID := strconv.Itoa(int(item.ID))

			assert.Equal(t, itemID, resp.AssignMeToItem.Id)
			assert.Equal(t, "100.00", resp.AssignMeToItem.Price)
			assert.Equal(t, *item.Name, resp.AssignMeToItem.Name)
			assert.Equal(t, userClaim.ID, resp.AssignMeToItem.SharedBy[0].Id)
			assert.Equal(t, user.Username, resp.AssignMeToItem.SharedBy[0].Username)
		}
	})

	t.Run("totalPayables of Me", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		// Creating 2 users
		user, err := app.Repositories.UserRepository.Create("john_doe", "password")
		if err != nil {
			t.Fatal(err)
		}

		user2, err := app.Repositories.UserRepository.Create("jane_doe", "password")
		if err != nil {
			t.Fatal(err)
		}

		// Receipt
		userID := strconv.Itoa(int(user.ID))
		receipt, _ := repository.NewReceipt(35000, "test receipt", userID)
		err = app.Repositories.ReceiptRepository.UnsafeCreate(&receipt)
		if err != nil {
			t.Fatal(err)
		}

		// Leaving this here for documentation on how to initialize a struct with embedded properties
		item := repository.Item{
			Items: g.Items{
				Name: repository.Nullable("Dumplings"), Price: int32(10001), ReceiptID: repository.BigInt(receipt.ID),
			},
		}

		err = app.Repositories.ItemRepository.Create(&item)
		if err != nil {
			slog.Error(err.Error())
			t.Fatal(err)
		}

		item2 := repository.Item{}
		item2.Name = repository.Nullable("Chicken")
		item2.ReceiptID = repository.BigInt(receipt.ID)
		item2.Price = int32(2788)

		err = app.Repositories.ItemRepository.Create(&item2)
		if err != nil {
			t.Fatal(err)
		}

		// Create user_orders
		// repository.UserOrdersRepository.
		_ = app.Repositories.UserOrdersRepository.Create(toString(user.ID), toString(item.ID))
		_ = app.Repositories.UserOrdersRepository.Create(toString(user2.ID), toString(item.ID))
		_ = app.Repositories.UserOrdersRepository.Create(toString(user.ID), toString(item2.ID))

		userClaim := auth.NewUserClaim(
			user.ID,
			user.Username,
			user.State,
		)
		accessToken, _ := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

		query := `query Me {
			me {
				totalPayables
			}
		}`

		var resp struct {
			Me struct {
				TotalPayables string
			}
		}

		option := func(bd *client.Request) {
			bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.Equal(t, "77.88", resp.Me.TotalPayables)
		}
	})

	t.Run("mutation DeleteMyReceipt", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		user, err := app.Repositories.UserRepository.Create("john_doe", "testing")
		if err != nil {
			t.Fatal(err)
		}

		receipt, _ := repository.NewReceipt(35000, "test receipt", toString(user.ID))
		err = app.Repositories.ReceiptRepository.UnsafeCreate(&receipt)
		if err != nil {
			t.Fatal(err)
		}

		item := repository.Item{} //{Items: model.Item{Name: "Dumplings", Price: 10000, ReceiptID: itemReceiptID}}
		item.Name = repository.Nullable("Dumplings")
		item.Price = int32(10000)
		item.ReceiptID = repository.BigInt(receipt.ID)

		err = app.Repositories.ItemRepository.Create(&item)
		if err != nil {
			t.Fatal(err)
		}

		userClaim := auth.NewUserClaim(
			user.ID,
			user.Username,
			user.State,
		)
		accessToken, _ := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

		query := fmt.Sprintf(`mutation DeleteMyReceipt {
			deleteMyReceipt(input: { id: "%d" })
		}`, receipt.ID)

		var resp struct {
			DeleteMyReceipt string
		}

		option := func(bd *client.Request) {
			bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.Equal(t, fmt.Sprintf("%d", receipt.ID), resp.DeleteMyReceipt)
		}
	})

	t.Run("query MyReceipts", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		receipt, err := CreateReceiptWithUser(app.DB, 9900, "receipt 1")
		if err != nil {
			slog.Error(err.Error())
			t.Fatal(err)
		}

		userClaim := auth.NewUserClaim(
			receipt.User.ID,
			receipt.User.Username,
			receipt.User.State,
		)
		accessToken, _ := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

		query := fmt.Sprintf(`query MyReceipts {
			myReceipts {
        description
        id
        userId
        total
        slug
      }
		}`)

		var resp struct {
			MyReceipts []*model.Receipt
		}

		option := func(bd *client.Request) {
			bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			if assert.Equal(t, 1, len(resp.MyReceipts)) {
				r := resp.MyReceipts[0]
				expectedID := strconv.Itoa(int(receipt.ID))
				expectedUserID := strconv.Itoa(int(receipt.User.ID))

				assert.Equal(t, "receipt 1", r.Description)
				assert.Equal(t, "99.00", r.Total)
				assert.Equal(t, expectedID, r.ID)
				assert.Equal(t, expectedUserID, r.UserID)
			}
		}
	})

	t.Run("mutation addItemToReceipt", func(t *testing.T) {
		defer TruncateAllTables(app.DB)

		// Creat User and their receipt
		user, _ := app.Repositories.UserRepository.Create("john_doe", "testing")
		userID := strconv.Itoa(int(user.ID))
		receipt, err := repository.NewReceipt(35000, "test receipt", userID)
		err = app.Repositories.ReceiptRepository.UnsafeCreate(&receipt)
		if err != nil {
			t.Fatal(err)
		}

		userClaim := auth.NewUserClaim(
			user.ID,
			user.Username,
			user.State,
		)
		accessToken, _ := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

		inputItem := struct {
			Name  string
			Price float32
		}{
			Name:  "Dumplings",
			Price: 18.00,
		}

		query := fmt.Sprintf(`mutation addItemToReceipt{
      addItemToReceipt(input: { receiptId: "%d", name: "%s", price: %f }) {
				name
				price
				id
				sharedBy {
					id
					username
				}
			}
		}`, receipt.ID, inputItem.Name, inputItem.Price)

		var resp struct {
			AddItemToReceipt struct {
				Price    string
				Name     string
				Id       string
				SharedBy []struct {
					Id       string
					Username string
				}
			}
		}

		option := func(bd *client.Request) {
			bd.HTTP.AddCookie(&http.Cookie{Name: string(internal.BearerTokenCookie), Value: accessToken})
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.NotNil(t, resp.AddItemToReceipt.Id)
			assert.Equal(t, "18.00", resp.AddItemToReceipt.Price)
			assert.Equal(t, inputItem.Name, resp.AddItemToReceipt.Name)

			assert.Empty(t, resp.AddItemToReceipt.SharedBy)
		}
	})
}
