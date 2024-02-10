package graph

import (
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/carlqt/ezsplit/graph/directive"
	"github.com/carlqt/ezsplit/internal"
	"github.com/stretchr/testify/assert"
)

func TestGraphqlServer(t *testing.T) {
	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDc1MTU3NTEsImlkIjoiMyIsInVzZXJuYW1lIjoidGVzdGVzdCJ9.Wak-YpGFivPsTKQjU3k0wf9HmX5qIo0w4hgYU6nLS_8"
	app := internal.NewApp()

	resolvers := &Resolver{Repositories: app.Repositories, Config: app.Config}
	config := Config{Resolvers: resolvers}

	config.Directives.Authenticated = directive.AuthDirective(app.Config.JWTSecret)
	srv := handler.NewDefaultServer(NewExecutableSchema(config))
	c := client.New(srv)

	t.Run("createUser mutation", func(t *testing.T) {

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

		client.AddHeader("Authorization", "Bearer "+accessToken)
		err := c.Post(query, &resp)

		assert.Nil(t, err)
		assert.Equal(t, "testest", resp.CreateUser.Username)
	})
}
