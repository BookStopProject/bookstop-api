package loader

import (
	"bookstop/app/book"
	"bookstop/app/inventory"
	"bookstop/app/location"
	"bookstop/app/user"
	"bookstop/app/userbook"
	"bookstop/graph/model"
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
				Fetch: func(keys []int) ([]*model.User, []error) {
					return user.LoadManyByIDs(origCtx, keys)
				},
			}),
			BookByID: *book.NewBookLoader(
				book.BookLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []string) ([]*book.Book, []error) {
						return book.LoadManyByIDs(origCtx, keys)
					},
				},
			),
			LocationByID: *location.NewLocationLoader(
				location.LocationLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*model.Location, []error) {
						return location.LoadManyByIDs(origCtx, keys)
					},
				},
			),
			UserBookByID: *userbook.NewUserBookLoader(
				userbook.UserBookLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*model.UserBook, []error) {
						return userbook.LoadManyByIDs(origCtx, keys)
					},
				},
			),
			InventoryByID: *inventory.NewInventoryLoader(
				inventory.InventoryLoaderConfig{
					Wait: Wait,
					Fetch: func(keys []int) ([]*model.Inventory, []error) {
						return inventory.LoadManyByIDs(origCtx, keys)
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
