package db

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func initDB() *pgx.Conn {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("No DATABASE_URL env")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Connected to database: %s", conn.Config().Database)

	return conn
}

var Conn = initDB()

func ParamRefsStr(l int) (paramrefs string) {
	for i := 1; i <= l; i++ {
		paramrefs += `$` + strconv.Itoa(i) + `,`
	}
	paramrefs = paramrefs[:len(paramrefs)-1] // remove last ","
	return
}
