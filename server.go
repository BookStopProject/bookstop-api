package main

import (
	"bookstop/auth"
	"bookstop/db"
	"bookstop/graph"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const defaultPort = "8080"
const gqlEndpoint = "/graphql"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	defer db.Conn.Close(context.Background())

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router := httprouter.New()
	router.Handler(http.MethodGet, "/playground", playground.Handler("GraphQL playground", gqlEndpoint))

	cor := cors.New(cors.Options{
		AllowOriginFunc: func(_ string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	})

	apiGraphQL := cor.Handler(auth.Middleware(srv))
	router.Handler(http.MethodGet, gqlEndpoint, apiGraphQL)
	router.Handler(http.MethodPost, gqlEndpoint, apiGraphQL)
	router.Handler(http.MethodOptions, gqlEndpoint, apiGraphQL)

	auth.Router(router)

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
