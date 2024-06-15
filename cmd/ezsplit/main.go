package main

import (
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carlqt/ezsplit/graph"
	"github.com/carlqt/ezsplit/graph/directive"
	"github.com/carlqt/ezsplit/internal"
)

const defaultPort = "8080"

func pong(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		slog.Error(err.Error())
	}
}

func main() {
	app := internal.NewApp()

	port := app.Config.Port
	if port == "" {
		port = defaultPort
	}

	c := graph.Config{Resolvers: &graph.Resolver{Repositories: app.Repositories, Config: app.Config}}
	c.Directives.Authenticated = directive.AuthDirective(app.Config.JWTSecret)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	// srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	// 	oc := graphql.GetOperationContext(ctx)
	// 	log.Printf("around: %s %s", oc.OperationName, oc.RawQuery)
	// 	return next(ctx)
	// })

	// TODO: Remove the playground handler in production.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	http.Handle("/query", internal.BearerTokenMiddleware(internal.InjectSetCookieMiddleware(srv)))
	http.Handle("/ping", http.HandlerFunc(pong))

	slog.Debug("connect to http://localhost:%s/ for GraphQL playground", "port", port)
	slog.Error(http.ListenAndServe(":"+port, nil).Error())
}
