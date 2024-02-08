package main

import (
	"log"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/carlqt/ezsplit/graph"
	"github.com/carlqt/ezsplit/graph/directive"
	"github.com/carlqt/ezsplit/internal"
)

func TestGETPlayers(t *testing.T) {
	app := internal.NewApp()

	resolvers := &graph.Resolver{Repositories: app.Repositories, Config: app.Config}
	config := graph.Config{Resolvers: resolvers}

	config.Directives.Authenticated = directive.AuthDirective(app.Config.JWTSecret)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(config))
	c := client.New(srv)

	t.Run("returns Pepper's score", func(t *testing.T) {

		query := "query { me { username } }"

		var resp struct {
			DirectiveArg *string
		}

		err := c.Post(query, &resp)
		if err != nil {
			t.Errorf(err.Error())
		}

		log.Println(resp)

		// response := httptest.NewRecorder()

		// got := response.Body.String()
		// want := "20"

		// if got != want {
		// 	t.Errorf("got %q, want %q", got, want)
		// }
	})
}
