package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.21 DO NOT EDIT.

import (
	"bookstop/models"
	"context"
	"strconv"
)

// Location is the resolver for the location field.
func (r *queryResolver) Location(ctx context.Context, id string) (*models.Location, error) {
	idNum, _ := strconv.Atoi(id)
	return models.FindLocationByID(ctx, idNum)
}

// Locations is the resolver for the locations field.
func (r *queryResolver) Locations(ctx context.Context) ([]*models.Location, error) {
	return models.FindAllLocations(ctx)
}
