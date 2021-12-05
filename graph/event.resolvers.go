package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/app/event"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/loader"
	"context"
	"strconv"
)

func (r *eventResolver) User(ctx context.Context, obj *model.Event) (*model.User, error) {
	intID, _ := strconv.Atoi(obj.UserID)
	return loader.For(ctx).UserByID.Load(intID)
}

func (r *queryResolver) Events(ctx context.Context) ([]*model.Event, error) {
	evts, err := event.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	results := make([]*model.Event, len(evts))
	for i, evt := range evts {
		results[i] = event.ToGraph(evt)
	}
	return results, nil
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

type eventResolver struct{ *Resolver }
