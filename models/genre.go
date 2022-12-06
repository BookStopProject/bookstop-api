package models

import (
	"bookstop/db"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Genre struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FindGenreByID(ctx context.Context, id int) (*Genre, error) {
	var genre Genre
	err := db.Conn.QueryRow(ctx, `SELECT
		id,
		name,
		description
	FROM
		public.genre
	WHERE
		id = $1
`, id).Scan(
		&genre.ID,
		&genre.Name,
		&genre.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &genre, nil
}

func FindAllGenres(ctx context.Context) ([]*Genre, error) {
	var genres []*Genre
	rows, err := db.Conn.Query(ctx, `SELECT
		id,
		name,
		description
	FROM
		public.genre
	ORDER BY
		name
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var genre Genre
		err = rows.Scan(
			&genre.ID,
			&genre.Name,
			&genre.Description)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &genre)
	}
	return genres, nil
}
