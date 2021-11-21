package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/book"
	"bookstop/browse"
	"bookstop/graph/model"
	"bookstop/loader"
	"context"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	bk, err := loader.For(ctx).BookByID.Load(id)
	if err != nil {
		return nil, err
	}
	return bk, nil
}

func (r *queryResolver) Books(ctx context.Context, ids []string) ([]*model.Book, error) {
	books, errs := loader.For(ctx).BookByID.LoadAll(ids)
	for i, err := range errs {
		if err != nil {
			graphql.AddError(ctx, gqlerror.Errorf("book "+strconv.Itoa(i)+": "+err.Error()))
		}
	}
	return books, nil
}

func (r *queryResolver) Browses(ctx context.Context) ([]*model.Browse, error) {
	now := time.Now()
	return browse.FindAll(ctx, &now)
}

func (r *queryResolver) Browse(ctx context.Context, id string) (*model.Browse, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return browse.FindByID(ctx, intID)
}

func (r *queryResolver) BrowseBooks(ctx context.Context, id string) ([]*model.Book, error) {
	intID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	books, _ := browse.FindBooksByBrowseID(ctx, intID)

	results := []*model.Book{}
	for _, book := range books {
		if book != nil {
			results = append(results, book)
		}
	}
	return results, nil
}

func (r *queryResolver) Search(ctx context.Context, query string, limit int, skip *int) ([]*model.Book, error) {
	startIndex := *skip
	return book.Search(ctx, query, limit, startIndex)
}
