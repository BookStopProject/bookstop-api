package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initDB() *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("No DATABASE_URL env")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Connected to database: %s", dbpool.Config().ConnConfig.Database)

	return dbpool
}

var Conn = initDB()

func ParamRefsStr(l int) (paramrefs string) {
	for i := 1; i <= l; i++ {
		paramrefs += `$` + strconv.Itoa(i) + `,`
	}
	paramrefs = paramrefs[:len(paramrefs)-1] // remove last ","
	return
}
