package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carlqt/ezsplit/graph"
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

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repositories: app.Repositories, Config: app.Config}}))

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
