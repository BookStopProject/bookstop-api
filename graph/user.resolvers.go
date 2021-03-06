package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/app/user"
	"bookstop/auth"
	"bookstop/graph/model"
	"bookstop/loader"
	"context"
	"strconv"
)

func (r *mutationResolver) MeUpdate(ctx context.Context, name *string, description *string) (*model.User, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	usr, err = user.UpdateByID(ctx, int(usr.ID.Int), name, description)
	if err != nil {
		return nil, err
	}
	return user.ToGraph(usr), nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	u := user.ToGraph(usr)
	if u != nil {
		creditInt := int(usr.Credit.Int)
		u.Credit = &creditInt
	}
	return u, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	intID, _ := strconv.Atoi(id)
	return loader.For(ctx).UserByID.Load(intID)
}
