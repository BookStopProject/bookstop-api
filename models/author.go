// CREATE TABLE public."author" (
// 	id serial PRIMARY KEY,
// 	name varchar(100) NOT NULL,
// 	description varchar(160),
// 	date_of_birth date varying(200) NOT NULL,
// 	date_of_death date varying(200) NOT NULL,
// );

package models

import (
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Author struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
	DateOfDeath *time.Time `json:"dateOfDeath"`
}

func FindAuthorByID(ctx context.Context, id int) (*Author, error) {
	var author Author
	err := db.Conn.QueryRow(ctx, `SELECT
		id,
		name,
		description,
		date_of_birth,
		date_of_death
	FROM
		public.author
	WHERE
		id = $1
`, id).Scan(
		&author.ID,
		&author.Name,
		&author.Description,
		&author.DateOfBirth,
		&author.DateOfDeath)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &author, nil
}

func SearchAuthor(ctx context.Context, name string) ([]Author, error) {
	var authors []Author
	rows, err := db.Conn.Query(ctx, `SELECT
		id,
		name,
		description,
		date_of_birth,
		date_of_death
	FROM
		public.author
	WHERE
		name ILIKE $1
`, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var author Author
		err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.Description,
			&author.DateOfBirth,
			&author.DateOfDeath)
		if err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}
