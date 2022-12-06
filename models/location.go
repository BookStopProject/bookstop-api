package models

import (
	"bookstop/db"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Location struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func FindLocationByID(ctx context.Context, id int) (*Location, error) {
	var location Location
	err := db.Conn.QueryRow(ctx, `SELECT
		id,
		name,
		address
	FROM
		public.location
	WHERE
		id = $1
	`, id).Scan(
		&location.ID,
		&location.Name,
		&location.Address,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &location, nil
}

func FindAllLocations(ctx context.Context) ([]*Location, error) {
	rows, err := db.Conn.Query(ctx, `SELECT
		id,
		name,
		address
	FROM
		public.location
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var locations []*Location
	for rows.Next() {
		var location Location
		err := rows.Scan(
			&location.ID,
			&location.Name,
			&location.Address,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}
	return locations, nil
}
