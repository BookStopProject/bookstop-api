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
	UserByID      user.UserLoader
	BookByID      book.BookLoader
	LocationByID  location.LocationLoader
	UserBookByID  userbook.UserBookLoader
	InventoryByID inventory.InventoryLoader
}

const Wait = 1 * time.Millisecond

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origCtx := r.Context()
		ctx := context.WithValue(origCtx, loadersKey, &Loaders{
			UserByID: *user.NewUserLoader(user.UserLoaderConfig{
				Wait: Wait,
				Fetch: func(keys []int) ([]*user.User, []error) {
					return user.FindManyByIDs(origCtx, keys)
				},
			}),
			BookByID: *book.NewBookLoader(
				book.BookLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []string) ([]*book.Book, []error) {
						return book.FindManyByIDs(origCtx, keys)
					},
				},
			),
			LocationByID: *location.NewLocationLoader(
				location.LocationLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*location.Location, []error) {
						return location.FindManyByIDs(origCtx, keys)
					},
				},
			),
			UserBookByID: *userbook.NewUserBookLoader(
				userbook.UserBookLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*userbook.UserBook, []error) {
						return userbook.FindManyByIDs(origCtx, keys)
					},
				},
			),
			InventoryByID: *inventory.NewInventoryLoader(
				inventory.InventoryLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*inventory.Inventory, []error) {
						return inventory.FindManyByIDs(origCtx, keys)
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
