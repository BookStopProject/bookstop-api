package models

import (
	"bookstop/db"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Book struct {
	ID             int     `json:"id"`
	Title          string  `json:"title"`
	Subtitle       *string `json:"subtitle"`
	ImageURL       *string `json:"imageUrl"`
	Description    *string `json:"description"`
	PublishedYear  int     `json:"publishedYear"`
	AuthorID       int     `json:"authorId"`
	GenreID        int     `json:"genreId"`
	TradeinCredit  int     `json:"tradeinCredit"`
	ExchangeCredit int     `json:"exchangeCredit"`
	Author         *Author `json:"author"`
	Genre          *Genre  `json:"genre"`
}

func FindBookByID(ctx context.Context, id int) (*Book, error) {
	var book Book
	book.Author = &Author{}
	book.Genre = &Genre{}
	err := db.Conn.QueryRow(ctx, `SELECT
	b.id,
	b.title,
	b.subtitle,
	b.description,
	b.image_url,
	b.published_year,
	b.tradein_credit,
	b.exchange_credit,
	a.id AS author_id,
	a.name AS author_name,
	g.id AS genre_id,
	g.name AS genre_name
FROM
	public.book b
	JOIN public.author a ON b.author_id = a.id
	JOIN public.genre g ON b.genre_id = g.id
WHERE
	b.id = $1
`, id).Scan(
		&book.ID,
		&book.Title,
		&book.Subtitle,
		&book.Description,
		&book.ImageURL,
		&book.PublishedYear,
		&book.TradeinCredit,
		&book.ExchangeCredit,
		&book.Author.ID,
		&book.Author.Name,
		&book.Genre.ID,
		&book.Genre.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func SearchBooks(ctx context.Context, query string, limit int, skip int) ([]*Book, error) {
	if limit == 0 {
		limit = 10
	}
	rows, err := db.Conn.Query(ctx, `SELECT
		b.id,
		b.title,
		b.subtitle,
		b.image_url,
		b.description,
		b.published_year,
		b.genre_id,
		b.tradein_credit,
		b.exchange_credit,
		a.id,
		a.name
	FROM
		public.book b
		JOIN public.author a ON b.author_id = a.id
	WHERE
		b.title ILIKE $1
		OR b.subtitle ILIKE $1
		OR a.name ILIKE $1
	LIMIT $2 OFFSET $3
`, "%"+query+"%", limit, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book

	for rows.Next() {
		var book Book
		book.Author = &Author{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Subtitle,
			&book.ImageURL,
			&book.Description,
			&book.PublishedYear,
			&book.GenreID,
			&book.TradeinCredit,
			&book.ExchangeCredit,
			&book.Author.ID,
			&book.Author.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	return books, nil
}

func FindBooksByAuthorId(ctx context.Context, authorId int) ([]Book, error) {
	rows, err := db.Conn.Query(ctx, `SELECT
		b.id,
		b.title,
		b.subtitle,
		b.image_url,
		b.description,
		b.published_year,
		b.genre_id,
		b.tradein_credit,
		b.exchange_credit,
		a.id,
		a.name
	FROM
		public.book b
		JOIN public.author a ON b.author_id = a.id
	WHERE
		author_id = $1
`, authorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		book.Author = &Author{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Subtitle,
			&book.ImageURL,
			&book.Description,
			&book.PublishedYear,
			&book.GenreID,
			&book.TradeinCredit,
			&book.ExchangeCredit,
			&book.Author.ID,
			&book.Author.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func FindBooksByGenreId(ctx context.Context, genreId int) ([]Book, error) {
	rows, err := db.Conn.Query(ctx, `SELECT
		b.id,
		b.title,
		b.subtitle,
		b.image_url,
		b.description,
		b.published_year,
		b.genre_id,
		b.tradein_credit,
		b.exchange_credit,
		a.id,
		a.name
	FROM
		public.book b
		JOIN public.author a ON b.author_id = a.id
	WHERE
		genre_id = $1
`, genreId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		book.Author = &Author{}
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Subtitle,
			&book.ImageURL,
			&book.Description,
			&book.PublishedYear,
			&book.GenreID,
			&book.TradeinCredit,
			&book.ExchangeCredit,
			&book.Author.ID,
			&book.Author.Name)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
