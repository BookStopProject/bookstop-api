package main

import (
	"bookstop/admin"
	"bookstop/auth"
	"bookstop/db"
	"bookstop/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

const gqlEndpoint = "/graphql"

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is not set")
	}

	defer db.Conn.Close()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	router := httprouter.New()
	router.Handler(http.MethodGet, "/playground", playground.Handler("GraphQL playground", gqlEndpoint))

	cor := cors.New(cors.Options{
		AllowOriginFunc: func(_ string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	apiGraphQL := auth.Middleware(srv)
	router.Handler(http.MethodGet, gqlEndpoint, apiGraphQL)
	router.Handler(http.MethodPost, gqlEndpoint, apiGraphQL)
	router.Handler(http.MethodOptions, gqlEndpoint, apiGraphQL)

	auth.Router(router)
	admin.Router(router)

	log.Printf("connect to http://localhost:%s/playground for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, cor.Handler(router)))
}
