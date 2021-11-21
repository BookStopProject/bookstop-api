package main

import (
	"bookstop/admin"
	"bookstop/auth"
	"bookstop/db"
	"bookstop/graph"
	"bookstop/graph/generated"
	"bookstop/loader"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

	defer db.Conn.Close()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router := httprouter.New()

	router.Handler(http.MethodGet, "/playground", playground.Handler("GraphQL playground", gqlEndpoint))

	cor := cors.New(cors.Options{
		AllowOriginFunc: func(_ string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	})

	apiGraphQL := cor.Handler(loader.Middleware(auth.Middleware(srv)))
	router.Handler(http.MethodGet, gqlEndpoint, apiGraphQL)
	router.Handler(http.MethodPost, gqlEndpoint, apiGraphQL)
	router.Handler(http.MethodOptions, gqlEndpoint, apiGraphQL)

	auth.Router(router)
	admin.Router(router)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
