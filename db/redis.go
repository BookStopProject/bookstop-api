package db

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func initRedisDB() *redis.Client {
	dbURL := os.Getenv("REDIS_URL")
	if dbURL == "" {
		log.Fatalln("No REDIS_URL env")
	}

	opt, err := redis.ParseURL(dbURL)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}

var RDB = initRedisDB()
