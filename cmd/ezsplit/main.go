package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carlqt/ezsplit/graph"
	"github.com/carlqt/ezsplit/graph/model"
	"github.com/carlqt/ezsplit/internal"
	"github.com/carlqt/ezsplit/internal/auth"
)

const defaultPort = "8080"

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	app := internal.NewApp()

	port := app.Config.Port
	if port == "" {
		port = defaultPort
	}

	c := graph.Config{Resolvers: &graph.Resolver{Repositories: app.Repositories, Config: app.Config}}
	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
		if role == model.RoleAuthenticatedUser {
			bearerToken := ctx.Value(auth.TokenKey).(string)

			claims, err := auth.ValidateBearerToken(bearerToken, app.Config.JWTSecret)
			if err != nil {
				log.Println(err)
				return nil, errors.New("unauthorized access")
			}

			ctx = context.WithValue(ctx, auth.UserClaimKey, claims)
		}

		return next(ctx)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	// srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	// 	oc := graphql.GetOperationContext(ctx)
	// 	log.Printf("around: %s %s", oc.OperationName, oc.RawQuery)
	// 	return next(ctx)
	// })

	// TODO: Remove the playground handler in production.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	http.Handle("/query", auth.BearerTokenMiddleware(srv))
	http.Handle("/ping", http.HandlerFunc(pong))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
