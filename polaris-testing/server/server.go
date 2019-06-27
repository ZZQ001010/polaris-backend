package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	coding_carefree "github.com/allstar/coding-carefree"
)

const defaultPort = "8080"

func main() {
	os.Setenv("UPPERIO_DB_DEBUG", "1")
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(coding_carefree.NewExecutableSchema(coding_carefree.Config{Resolvers: &coding_carefree.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
