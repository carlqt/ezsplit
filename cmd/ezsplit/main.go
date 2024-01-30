package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/carlqt/ezsplit/graph"
	"github.com/carlqt/ezsplit/internal"
	middleware "github.com/carlqt/ezsplit/internal/middlewares"
)

const defaultPort = "8080"

func main() {
	app := internal.NewApp()

	port := app.Config.Port
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Repositories: app.Repositories}}))

	// srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	// 	oc := graphql.GetOperationContext(ctx)
	// 	log.Printf("around: %s %s", oc.OperationName, oc.RawQuery)
	// 	return next(ctx)
	// })

	// TODO: Remove the playground handler in production.
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// http.Handle("/query", handlers.CombinedLoggingHandler(os.Stdout, srv))

	// TODO: Add authentication middleware
	http.Handle("/query", middleware.AuthMiddleware(srv, *app.Config))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
