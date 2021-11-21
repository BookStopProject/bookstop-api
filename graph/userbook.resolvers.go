package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/auth"
	"bookstop/graph/generated"
	"bookstop/graph/model"
	"bookstop/loader"
	"bookstop/user"
	"bookstop/userbook"
	"context"
	"errors"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) UserBookAdd(ctx context.Context, bookID string, startedAt *string, endedAt *string) (*model.UserBook, error) {
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	ub, err := userbook.Create(ctx, int(usr.ID.Int), bookID, startedAt, endedAt)
	if err != nil {
		return nil, err
	}
	return userbook.ToGraph(ub), nil
}

func (r *mutationResolver) UserBookEdit(ctx context.Context, id string, startedAt *string, endedAt *string) (*model.UserBook, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, auth.ErrUnauthorized
	}
	owned, err := userbook.IsOwner(ctx, int(usr.ID.Int), intID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, auth.ErrForbidden
	}
	ub, err := userbook.UpdateByID(ctx, intID, startedAt, endedAt)
	if err != nil {
		return nil, err
	}
	return userbook.ToGraph(ub), nil
}

func (r *mutationResolver) UserBookDelete(ctx context.Context, id string) (bool, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return false, err
	}
	usr, err := auth.ForContext(ctx)
	if err != nil {
		return false, err
	}
	if usr == nil {
		return false, auth.ErrUnauthorized
	}
	owned, err := userbook.IsOwner(ctx, int(usr.ID.Int), intID)
	if err != nil {
		return false, err
	}
	if !owned {
		return false, auth.ErrForbidden
	}
	return userbook.DeleteByID(ctx, intID)
}

func (r *queryResolver) UserBook(ctx context.Context, id string) (*model.UserBook, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	userBook, err := loader.For(ctx).UserBookByID.Load(intID)
	if err != nil {
		return nil, err
	}
	return userbook.ToGraph(userBook), nil
}

func (r *queryResolver) UserBooks(ctx context.Context, userID *string, mine *bool) ([]*model.UserBook, error) {
	var results []*model.UserBook
	var userBooks []*userbook.UserBook
	var errs []error
	if userID != nil {
		intUserID, err := strconv.Atoi(*userID)
		if err != nil {
			return nil, err
		}
		userBooks, errs = userbook.FindManyByUserID(ctx, intUserID)
	} else if *mine {
		usr, err := auth.ForContext(ctx)
		if err != nil {
			return nil, err
		}
		if usr == nil {
			return nil, auth.ErrUnauthorized
		}
		userBooks, errs = userbook.FindManyByUserID(ctx, int(usr.ID.Int))
	} else {
		return nil, errors.New("must provide either userID or mine = true")
	}
	for i, ub := range userBooks {
		if errs[i] != nil {
			graphql.AddError(ctx, gqlerror.Errorf("userbook "+strconv.Itoa(i)+": "+errs[i].Error()))
		}
		results = append(results, userbook.ToGraph(ub))
	}
	return results, nil
}

func (r *userBookResolver) Book(ctx context.Context, obj *model.UserBook) (*model.Book, error) {
	return loader.For(ctx).BookByID.Load(obj.BookID)
}

func (r *userBookResolver) User(ctx context.Context, obj *model.UserBook) (*model.User, error) {
	intID, _ := strconv.Atoi(obj.UserID)
	usr, err := loader.For(ctx).UserByID.Load(intID)
	if err != nil {
		return nil, err
	}
	return user.ToGraph(usr), nil
}

// UserBook returns generated.UserBookResolver implementation.
func (r *Resolver) UserBook() generated.UserBookResolver { return &userBookResolver{r} }

type userBookResolver struct{ *Resolver }
