package loader

import (
	"bookstop/book"
	"bookstop/inventory"
	"bookstop/location"
	"bookstop/user"
	"bookstop/userbook"
	"context"
	"net/http"
	"time"
)

const loadersKey = "dataloaders"

type Loaders struct {
	UserById      user.UserLoader
	BookById      book.BookLoader
	LocationById  location.LocationLoader
	UserBookById  userbook.UserBookLoader
	InventoryById inventory.InventoryLoader
}

const Wait = 1 * time.Millisecond

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origCtx := r.Context()
		ctx := context.WithValue(origCtx, loadersKey, &Loaders{
			UserById: *user.NewUserLoader(user.UserLoaderConfig{
				Wait: Wait,
				Fetch: func(keys []int) ([]*user.User, []error) {
					return user.FindManyByIds(origCtx, keys)
				},
			}),
			BookById: *book.NewBookLoader(
				book.BookLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []string) ([]*book.Book, []error) {
						return book.FindManyByIds(origCtx, keys)
					},
				},
			),
			LocationById: *location.NewLocationLoader(
				location.LocationLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*location.Location, []error) {
						return location.FindManyByIds(origCtx, keys)
					},
				},
			),
			UserBookById: *userbook.NewUserBookLoader(
				userbook.UserBookLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*userbook.UserBook, []error) {
						return userbook.FindManyByIds(origCtx, keys)
					},
				},
			),
			InventoryById: *inventory.NewInventoryLoader(
				inventory.InventoryLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*inventory.Inventory, []error) {
						return inventory.FindManyByIds(origCtx, keys)
					},
				},
			),
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
