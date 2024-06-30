package integration_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
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
	defer TruncateAllTables(app.DB)

	t.Run("loginUser mutation", func(t *testing.T) {
		t.Run("when there are no inputs", func(t *testing.T) {
			defer TruncateAllTables(app.DB)

			_, err := CreateUser(app.DB, "mutation_user160")
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

			_, err := CreateUser(app.DB, "mutation_user160")
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

			if assert.Nil(t, err) {
				assert.Equal(t, "mutation_user160", resp.LoginUser.Username)
			}
		})

		t.Run("when password is wrong", func(t *testing.T) {
			defer TruncateAllTables(app.DB)

			_, err := CreateUser(app.DB, "mutation_user160")
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

		user, err := CreateUser(app.DB, "fake_user")
		if err != nil {
			t.Fatal(err)
		}

		t.Run("with Receipts field", func(t *testing.T) {
			userClaim := auth.UserClaim{
				ID:       user.ID,
				Username: user.Username,
			}
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

				receipt := repository.Receipt{Description: "test receipt", Total: 35000, UserID: user.ID}
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

					assert.Equal(t, user.Username, resp.Me.Username)
					assert.Equal(t, receipt.ID, responseReceipt.ID)
					assert.Equal(t, receipt.Description, responseReceipt.Description)

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

					assert.Equal(t, user.ID, resp.Me.Id)
					assert.Empty(t, responseReceipt)
				}
			})
		})

		t.Run("when jwt exists", func(t *testing.T) {
			userClaim := auth.UserClaim{
				ID:       user.ID,
				Username: user.Username,
			}
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
				assert.Equal(t, user.ID, resp.Me.Id)
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

		user, err := CreateUser(app.DB, "fake_user")
		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := user.getAuthToken(app.Config.JWTSecret)
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

		receipt := repository.Receipt{Description: "test receipt", Total: 35000, UserID: user.ID}
		err = app.Repositories.ReceiptRepository.Create(&receipt)
		if err != nil {
			t.Fatal(err)
		}

		item := repository.Item{Name: "Dumplings", Price: 10000, ReceiptID: receipt.ID}
		err = app.Repositories.ItemRepository.Create(&item)
		if err != nil {
			t.Fatal(err)
		}

		userClaim := auth.UserClaim{
			ID:       user.ID,
			Username: user.Username,
		}
		accessToken, _ := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

		query := fmt.Sprintf(`mutation assignMeToItem{
			assignMeToItem(input: { itemId: "%s" }) {
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
			assert.Equal(t, item.Name, resp.AssignMeToItem.Name)
			assert.Equal(t, user.ID, resp.AssignMeToItem.SharedBy[0].Id)
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
		receipt := repository.Receipt{Description: "test receipt", Total: 35000, UserID: user.ID}
		err = app.Repositories.ReceiptRepository.Create(&receipt)
		if err != nil {
			t.Fatal(err)
		}

		item := repository.Item{Name: "Dumplings", Price: 10000, ReceiptID: receipt.ID}
		err = app.Repositories.ItemRepository.Create(&item)
		if err != nil {
			t.Fatal(err)
		}

		item2 := repository.Item{Name: "Chicken", Price: 2788, ReceiptID: receipt.ID}
		err = app.Repositories.ItemRepository.Create(&item2)
		if err != nil {
			t.Fatal(err)
		}

		// Create user_orders
		// repository.UserOrdersRepository.
		_ = app.Repositories.UserOrdersRepository.Create(user.ID, item.ID)
		_ = app.Repositories.UserOrdersRepository.Create(user2.ID, item.ID)
		_ = app.Repositories.UserOrdersRepository.Create(user.ID, item2.ID)

		userClaim := auth.UserClaim{
			ID:       user.ID,
			Username: user.Username,
		}
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
}
