package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/app/location"
	"bookstop/graph/model"
	"context"
)

func (r *queryResolver) Locations(ctx context.Context) ([]*model.Location, error) {
	return location.FindAllLocations(ctx)
}
