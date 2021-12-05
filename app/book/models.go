package book

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/go-redis/redis/v8"
	"google.golang.org/api/books/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type Book = model.Book

func NewBookService(ctx context.Context) (*books.Service, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatalln("No GOOGLE_API_KEY env")
	}
	return books.NewService(ctx, option.WithAPIKey(apiKey))
}

var apiUsedFields = []googleapi.Field{"id", "volumeInfo(title,authors,description,publishedDate,subtitle,imageLinks,industryIdentifiers)"}
var searchApiUsedFields = []googleapi.Field{"items/" + apiUsedFields[0], "items/" + apiUsedFields[1]}

const redisBookPrefix = "book:"
const redisCacheExpiration = 14 * 24 * time.Hour // 14 days

func setCache(ctx context.Context, books []*model.Book) error {
	_, err := db.RDB.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, book := range books {
			if book != nil {
				b, err := json.Marshal(book)
				if err == nil {
					pipe.Set(ctx, redisBookPrefix+book.ID, string(b), redisCacheExpiration)
				}
			}
		}
		return nil
	})
	return err
}

func getCache(ctx context.Context, ids []string) ([]*model.Book, error) {
	keys := make([]string, len(ids))
	for idx, id := range ids {
		keys[idx] = redisBookPrefix + id
	}
	strBooks, err := db.RDB.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make([]*model.Book, len(ids))

	for idx, strBook := range strBooks {
		if strBook != nil {
			var book *model.Book = &model.Book{}
			err := json.Unmarshal([]byte(strBook.(string)), book)
			if err == nil {
				results[idx] = book
			}
		}
	}

	return results, nil
}

func findByIDViaAPI(ctx context.Context, id string) (*model.Book, error) {
	srv, err := NewBookService(ctx)

	if err != nil {
		return nil, err
	}
	volume, err := srv.Volumes.Get(id).Fields(apiUsedFields...).Do()
	if err != nil {
		if err.(*googleapi.Error).Code == 404 {
			return nil, nil
		}
		return nil, err
	}

	if volume == nil {
		return nil, nil
	}

	return ToGraph(toMd(volume)), nil
}

func FindByID(ctx context.Context, id string) (*model.Book, error) {
	cached, err := getCache(ctx, []string{id})
	if err != nil {
		return nil, err
	}
	if cached[0] != nil {
		return cached[0], nil
	}

	book, err := findByIDViaAPI(ctx, id)

	if err != nil {
		return nil, err
	}

	setCache(ctx, []*model.Book{book})

	return book, nil
}

func Search(ctx context.Context, query string, limit int, startIndex int) (books []*model.Book, err error) {
	srv, err := NewBookService(ctx)
	if err != nil {
		return nil, err
	}
	volumes, err := srv.Volumes.List(query).Fields(searchApiUsedFields...).MaxResults(int64(limit)).StartIndex(int64(startIndex)).Do()
	if err != nil {
		return nil, err
	}
	for _, volume := range volumes.Items {
		books = append(books, ToGraph(toMd(volume)))
	}
	return
}
