package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/app/thought"
	"bookstop/auth"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/loader"
	"context"
	"strconv"
)

func (r *mutationResolver) ThoughtCreate(ctx context.Context, text string, bookID *string) (*model.Thought, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	t, err := thought.Create(ctx, int(usr.ID.Int), text, bookID)
	if err != nil {
		return nil, err
	}
	return thought.ToGraph(t), nil
}

func (r *mutationResolver) ThoughtDelete(ctx context.Context, id string) (bool, error) {
	intID, _ := strconv.Atoi(id)
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return false, err
	}
	if usr == nil {
		return false, auth.ErrUnauthorized
	}
	owned, err := thought.IsOwner(ctx, int(usr.ID.Int), intID)
	if err != nil {
		return false, err
	}
	if !owned {
		return false, auth.ErrForbidden
	}
	return thought.DeleteByID(ctx, intID)
}

func (r *queryResolver) Thoughts(ctx context.Context, userID *string, limit int, before *int) ([]*model.Thought, error) {
	var results []*thought.Thought
	var err error
	if err != nil {
		return nil, err
	}
	if userID != nil {
		intID, errConv := strconv.Atoi(*userID)
		if errConv != nil {
			return nil, errConv
		}
		results, err = thought.FindManyByUserID(ctx, intID, limit, before)
	} else {
		results, err = thought.FindAll(ctx, limit, before)
	}
	if err != nil {
		return nil, err
	}
	thts := make([]*model.Thought, len(results))
	for i, t := range results {
		thts[i] = thought.ToGraph(t)
	}
	return thts, nil
}

func (r *thoughtResolver) User(ctx context.Context, obj *model.Thought) (*model.User, error) {
	intID, _ := strconv.Atoi(obj.UserID)
	return loader.For(ctx).UserByID.Load(intID)
}

func (r *thoughtResolver) Book(ctx context.Context, obj *model.Thought) (*model.Book, error) {
	if obj.BookID == nil {
		return nil, nil
	}
	return loader.For(ctx).BookByID.Load(*obj.BookID)
}

// Thought returns generated.ThoughtResolver implementation.
func (r *Resolver) Thought() generated.ThoughtResolver { return &thoughtResolver{r} }

type thoughtResolver struct{ *Resolver }
