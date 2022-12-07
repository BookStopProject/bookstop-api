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

func (p *Post) IsOwner(ctx context.Context, userID int) bool {
	return p.UserID == userID
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
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.post
		JOIN public.book ON post.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.user ON post.user_id = "user".id
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

func FindAllPosts(ctx context.Context, limit int, before *int) ([]*Post, error) {
	if limit == 0 {
		limit = 10
	}

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
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.post
		JOIN public."book" ON post.book_id = book.id
		JOIN public."author" ON book.author_id = author.id
		JOIN public."user" ON post.user_id = "user".id
	WHERE post.id < $2
	LIMIT $1
	`, limit, before)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func FindPostsByUserID(ctx context.Context, userID int, limit int, before *int) ([]*Post, error) {
	if limit == 0 {
		limit = 10
	}

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
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.post
		JOIN public.book ON post.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.user ON post.user_id = "user".id
	WHERE
		post.user_id = $1
		AND post.id < $3
	LIMIT $2
	`, userID, limit, before)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func FindPostsByBookID(ctx context.Context, bookID int, limit int, before *int) ([]*Post, error) {
	if limit == 0 {
		limit = 10
	}

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
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.post
		JOIN public.book ON post.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.user ON post.user_id = "user".id
	WHERE
		post.book_id = $1
		AND post.id < $3
	LIMIT $2
	`, bookID, limit, before)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

func CreatePost(ctx context.Context, post *Post) (*Post, error) {
	// insert post and return the created post
	err := db.Conn.QueryRow(ctx, `INSERT INTO public.post (text, creation_time, book_id, user_id, is_recommending)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, post.Text, post.CreationTime, post.BookID, post.UserID, post.IsRecommending).Scan(&post.ID)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func UpdatePost(ctx context.Context, post *Post) (*Post, error) {
	// update post and return the updated post
	_, err := db.Conn.Exec(ctx, `UPDATE public.post SET text = $1, creation_time = $2, book_id = $3, user_id = $4, is_recommending = $5
	WHERE id = $6`, post.Text, post.CreationTime, post.BookID, post.UserID, post.IsRecommending, post.ID)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func DeletePost(ctx context.Context, postID int) error {
	// delete post
	_, err := db.Conn.Exec(ctx, `DELETE FROM public.post WHERE id = $1`, postID)

	if err != nil {
		return err
	}

	return nil
}
