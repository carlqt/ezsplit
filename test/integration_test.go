package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/carlqt/ezsplit/graph"
	"github.com/carlqt/ezsplit/graph/directive"
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
	c := client.New(srv)

	t.Run("createUser mutation", func(t *testing.T) {
		defer truncateAllTables(app.DB)

		query := `mutation createUser {
			createUser(input: {username: "testest" }) {
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
			assert.Equal(t, "testest", resp.CreateUser.Username)
		}
	})

	t.Run("query Me", func(t *testing.T) {
		defer truncateAllTables(app.DB)

		user, err := CreateUser(app.DB, "fake_user")
		if err != nil {
			t.Fatal(err)
		}

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
			ctx := context.WithValue(context.Background(), auth.TokenKey, accessToken)
			bd.HTTP = bd.HTTP.WithContext(ctx)
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.Equal(t, user.Username, resp.Me.Username)
			assert.Equal(t, user.ID, resp.Me.Id)
		}

	})

	t.Run("mutation createMyReceipt", func(t *testing.T) {
		defer truncateAllTables(app.DB)

		user, err := CreateUser(app.DB, "fake_user")
		if err != nil {
			t.Fatal(err)
		}

		accessToken, err := user.getAuthToken(app.Config.JWTSecret)
		if err != nil {
			t.Fatal(err)
		}

		query := `mutation createMyReceipt {
			createMyReceipt(input: {description: "test receipt", price: 350 }) {
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
			ctx := context.WithValue(context.Background(), auth.TokenKey, accessToken)
			bd.HTTP = bd.HTTP.WithContext(ctx)
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.Equal(t, "test receipt", resp.CreateMyReceipt.Description)
			assert.Equal(t, "350.00", resp.CreateMyReceipt.Total)
		}
	})

	t.Run("mutation assignMeToItem", func(t *testing.T) {
		defer truncateAllTables(app.DB)

		user, err := app.Repositories.UserRepository.Create("john_doe")
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
			ctx := context.WithValue(context.Background(), auth.TokenKey, accessToken)
			bd.HTTP = bd.HTTP.WithContext(ctx)
		}

		err = c.Post(query, &resp, option)

		if assert.Nil(t, err) {
			assert.Equal(t, item.Name, resp.AssignMeToItem.Name)
			assert.Equal(t, user.ID, resp.AssignMeToItem.SharedBy[0].Id)
			assert.Equal(t, user.Username, resp.AssignMeToItem.SharedBy[0].Username)
		}
	})
}

func truncateAllTables(db *sql.DB) {
	query := `
		DO $$ DECLARE
			r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname =current_schema()) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`

	slog.Debug("Clearing data")

	_, err := db.Exec(query)
	if err != nil {
		slog.Error(err.Error())
	}
}
