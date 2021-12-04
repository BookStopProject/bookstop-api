package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/exchange"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/loader"
	"context"
	"strconv"
)

func (r *exchangeResolver) UserBookOld(ctx context.Context, obj *model.Exchange) (*model.UserBook, error) {
	intID, _ := strconv.Atoi(obj.UserBookIDOld)
	return loader.For(ctx).UserBookByID.Load(intID)
}

func (r *exchangeResolver) UserBookNew(ctx context.Context, obj *model.Exchange) (*model.UserBook, error) {
	intID, _ := strconv.Atoi(obj.UserBookIDNew)
	return loader.For(ctx).UserBookByID.Load(intID)
}

func (r *queryResolver) Exchanges(ctx context.Context, userBookID string) ([]*model.Exchange, error) {
	intID, _ := strconv.Atoi(userBookID)
	return exchange.FindByUserBookID(ctx, intID)
}

// Exchange returns generated.ExchangeResolver implementation.
func (r *Resolver) Exchange() generated.ExchangeResolver { return &exchangeResolver{r} }

type exchangeResolver struct{ *Resolver }
