package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/exchange"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/loader"
	"bookstop/userbook"
	"context"
	"strconv"
)

func (r *exchangeResolver) UserBookOld(ctx context.Context, obj *model.Exchange) (*model.UserBook, error) {
	intID, _ := strconv.Atoi(obj.UserBookIDOld)
	ub, err := loader.For(ctx).UserBookByID.Load(intID)
	if err != nil {
		return nil, err
	}
	return userbook.ToGraph(ub), nil
}

func (r *exchangeResolver) UserBookNew(ctx context.Context, obj *model.Exchange) (*model.UserBook, error) {
	intID, _ := strconv.Atoi(obj.UserBookIDNew)
	ub, err := loader.For(ctx).UserBookByID.Load(intID)
	if err != nil {
		return nil, err
	}
	return userbook.ToGraph(ub), nil
}

func (r *queryResolver) Exchanges(ctx context.Context, userBookID string) ([]*model.Exchange, error) {
	intID, err := strconv.Atoi(userBookID)
	if err != nil {
		return nil, err
	}
	return exchange.FindExchangesByUserBookID(ctx, intID)
}

// Exchange returns generated.ExchangeResolver implementation.
func (r *Resolver) Exchange() generated.ExchangeResolver { return &exchangeResolver{r} }

type exchangeResolver struct{ *Resolver }
