package graph

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/carlqt/ezsplit/graph/directive"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
	"github.com/carlqt/ezsplit/internal/repository"
	"github.com/stretchr/testify/assert"
)

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Reset = "\033[0m"
)

func printRed(s string) {
	fmt.Printf("%s%s%s", Red, s, Reset)
}

func printGreen(s string) {
	fmt.Printf("%s%s%s", Green, s, Reset)
}

func printBlue(s string) {
	fmt.Printf("%s%s%s", Blue, s, Reset)
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	printBlue(">> Setup Test\n")

	return func(tb testing.TB) {

		printBlue(">> Teardown Test\n")
	}
}

func createUser(t *testing.T, db *sql.DB, user *repository.User) {
	user.Username = "john_watson_test"
	err := db.QueryRow("INSERT INTO users (username) VALUES ($1) RETURNING id", user.Username).Scan(&user.ID)

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		log.Println("Cleaning up user=" + user.ID)
		result, err := db.Exec("DELETE FROM users WHERE id = $1", user.ID)

		if err == nil {
			count, err := result.RowsAffected()
			if err == nil {
				log.Println("Deleted rows: ", count)
			} else {
				log.Println(err)
			}
		}
	})
}

func TestGraphqlServer(t *testing.T) {
	// Setting up the server
	app := internal.NewApp()
	resolvers := &Resolver{Repositories: app.Repositories, Config: app.Config}
	config := Config{Resolvers: resolvers}
	config.Directives.Authenticated = directive.AuthDirective(app.Config.JWTSecret)
	srv := handler.NewDefaultServer(NewExecutableSchema(config))
	c := client.New(srv)

	t.Run("createUser mutation", func(t *testing.T) {
		teardown := setupTest(t)
		defer teardown(t)

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

		// client.AddHeader("Authorization", "Bearer "+accessToken)
		err := c.Post(query, &resp)

		assert.Nil(t, err)
		assert.Equal(t, "testest", resp.CreateUser.Username)
	})

	t.Run("query Me", func(t *testing.T) {
		var user repository.User
		createUser(t, app.DB, &user)

		userClaim := auth.UserClaim{
			ID:       user.ID,
			Username: user.Username,
		}
		accessToken, err := auth.CreateAndSignToken(userClaim, app.Config.JWTSecret)

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
		if err != nil {
			slog.Error(err.Error())
		}

		assert.Equal(t, user.Username, resp.Me.Username)
		assert.Equal(t, user.ID, resp.Me.Id)
	})
}

func addContext(token string) client.Option {
	return func(bd *client.Request) {
		ctx := context.WithValue(context.Background(), auth.TokenKey, token)
		bd.HTTP = bd.HTTP.WithContext(ctx)
		log.Println("Added Context")
	}
}
