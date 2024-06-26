package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sigmawq/grpc-service/api/graph"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serviceHost := os.Getenv("SERVICE_HOST")
	if serviceHost == "" {
		serviceHost = "localhost:9000"
	}

	log.Printf("PORT=%v, SERVICE_HOST=%v", port, serviceHost)

	err := graph.InitializeGraphQLClient(serviceHost)
	if err != nil {
		os.Exit(1)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
