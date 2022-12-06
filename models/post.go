package models

import (
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Post struct {
	ID             int       `json:"id"`
	Text           string    `json:"text"`
	CreationTime   time.Time `json:"creationTime"`
	BookID         int       `json:"bookId"`
	UserID         int       `json:"userId"`
	IsRecommending bool      `json:"isRecommending"`
	Book           *Book     `json:"book"`
	User           *User     `json:"user"`
}

func FindPostByID(ctx context.Context, id int) (*Post, error) {
	var post Post
	post.Book = &Book{}
	post.Book.Author = &Author{}
	post.User = &User{}
	err := db.Conn.QueryRow(ctx, `SELECT
		post.id,
		post.text,
		post.creation_time,
		post.book_id,
		post.user_id,
		post.is_recommending,
		book.id,
		book.title,
		book.subtitle,
		book.author_id,
		author.id,
		author.name,
		user.id,
		user.name,
		user.profile_picture
	FROM
		public.post
		JOIN public.book ON post.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.user ON post.user_id = user.id
	WHERE
		post.id = $1
	`, id).Scan(
		&post.ID,
		&post.Text,
		&post.CreationTime,
		&post.BookID,
		&post.UserID,
		&post.IsRecommending,
		&post.Book.ID,
		&post.Book.Title,
		&post.Book.Subtitle,
		&post.Book.AuthorID,
		&post.Book.Author.ID,
		&post.Book.Author.Name,
		&post.User.ID,
		&post.User.Name,
		&post.User.ProfilePicture,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

func FindPostsByUserID(ctx context.Context, userID int) ([]*Post, error) {
	rows, err := db.Conn.Query(ctx, `SELECT
		post.id,
		post.text,
		post.creation_time,
		post.book_id,
		post.user_id,
		post.is_recommending,
		book.id,
		book.title,
		book.subtitle,
		book.author_id,
		author.id,
		author.name,
		user.id,
		user.name,
		user.profile_picture
	FROM
		public.post
		JOIN public.book ON post.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.user ON post.user_id = user.id
	WHERE
		post.user_id = $1
	`, userID)

	if err != nil {
		return nil, err
	}

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		post.Book = &Book{}
		post.Book.Author = &Author{}
		post.User = &User{}
		err := rows.Scan(
			&post.ID,
			&post.Text,
			&post.CreationTime,
			&post.BookID,
			&post.UserID,
			&post.IsRecommending,
			&post.Book.ID,
			&post.Book.Title,
			&post.Book.Subtitle,
			&post.Book.AuthorID,
			&post.Book.Author.ID,
			&post.Book.Author.Name,
			&post.User.ID,
			&post.User.Name,
			&post.User.ProfilePicture,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func FindPostsByBookID(ctx context.Context, bookID int) ([]*Post, error) {
	rows, err := db.Conn.Query(ctx, `SELECT
		post.id,
		post.text,
		post.creation_time,
		post.book_id,
		post.user_id,
		post.is_recommending,
		book.id,
		book.title,
		book.subtitle,
		book.author_id,
		author.id,
		author.name,
		user.id,
		user.name,
		user.profile_picture
	FROM
		public.post
		JOIN public.book ON post.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.user ON post.user_id = user.id
	WHERE
		post.book_id = $1
	`, bookID)

	if err != nil {
		return nil, err
	}

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		post.Book = &Book{}
		post.Book.Author = &Author{}
		post.User = &User{}
		err := rows.Scan(
			&post.ID,
			&post.Text,
			&post.CreationTime,
			&post.BookID,
			&post.UserID,
			&post.IsRecommending,
			&post.Book.ID,
			&post.Book.Title,
			&post.Book.Subtitle,
			&post.Book.AuthorID,
			&post.Book.Author.ID,
			&post.Book.Author.Name,
			&post.User.ID,
			&post.User.Name,
			&post.User.ProfilePicture,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
