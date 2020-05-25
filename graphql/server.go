package graphql

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func Serve() error {
	// GraphQL Server
	graphqlPort := ":3000"
	if os.Getenv("GRAPHQL_PORT") != "" {
		graphqlPort = os.Getenv("GRAPHQL_PORT")
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", handler.NewDefaultServer(NewExecutableSchema(
		Config{Resolvers: &Resolver{}}),
	))

	graphqlServer := &http.Server{
		Addr:    graphqlPort,
		Handler: mux,
	}
	defer func() {
		if err := graphqlServer.Shutdown(context.Background()); err != nil {
			log.Panicf("unable to shutdown graphql server gracefully: %+v", err)
		}
	}()

	log.Printf("GraphQL Playground is available at http://localhost%s/graphql/playground", graphqlPort)
	return graphqlServer.ListenAndServe()
}
