package main

import (
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
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	defer db.Conn.Close()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))

	cor := cors.New(cors.Options{
		AllowOriginFunc: func(_ string) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"*"},
	})

	http.Handle("/graphql", cor.Handler(loader.Middleware(auth.Middleware(srv))))

	http.HandleFunc("/auth", auth.RedirectHandler)
	http.HandleFunc("/auth/callback", auth.CallbackHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
