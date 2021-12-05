package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
)

func (r *queryResolver) HomeStats(ctx context.Context) (*model.HomeStats, error) {
	stats := model.HomeStats{}

	db.Conn.QueryRow(ctx, "SELECT COUNT(*) FROM public.user").Scan(&stats.UserCount)
	db.Conn.QueryRow(ctx, "SELECT COUNT(*) FROM public.inventory").Scan(&stats.ExchangeCount)
	db.Conn.QueryRow(ctx, "SELECT COUNT(*) FROM public.thought").Scan(&stats.PostCount)

	return &stats, nil
}
