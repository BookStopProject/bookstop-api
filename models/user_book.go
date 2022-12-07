package models

import (
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserBook struct {
	ID         int  `json:"id"`
	UserID     int  `json:"user_id"`
	BookID     int  `json:"book_id"`
	BookCopyID *int `json:"book_copy_id"`
	// Date user starts reading the book
	StartDate *time.Time `json:"startDate"`
	// Date user finishes reading the book
	EndDate *time.Time `json:"endDate"`
	Book    *Book      `json:"book"`
	User    *User      `json:"user"`
}

func (ub *UserBook) IsOwner(userID int) bool {
	return ub.UserID == userID
}

func FindUserBookByID(ctx context.Context, id int) (*UserBook, error) {
	// Find user book by id and join with book and user
	var userBook UserBook
	userBook.Book = &Book{}
	userBook.Book.Author = &Author{}
	userBook.Book.Genre = &Genre{}
	userBook.User = &User{}
	err := db.Conn.QueryRow(ctx, `SELECT
		user_book.id,
		user_book.user_id,	
		user_book.book_id,
		user_book.start_date,
		user_book.end_date,
		user_book.book_copy_id,
		book.id,
		book.title,
		book.subtitle,
		book.description,
		book.published_year,
		book.author_id,
		book.genre_id,
		author.id,
		author.name,
		genre.id,
		genre.name,
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.user_book
		JOIN public.book ON user_book.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.genre ON book.genre_id = genre.id
		JOIN public.user ON user_book.user_id = user.id
	WHERE
		user_book.id = $1
	`, id).Scan(
		&userBook.ID,
		&userBook.UserID,
		&userBook.BookID,
		&userBook.StartDate,
		&userBook.EndDate,
		&userBook.BookCopyID,
		&userBook.Book.ID,
		&userBook.Book.Title,
		&userBook.Book.Subtitle,
		&userBook.Book.Description,
		&userBook.Book.PublishedYear,
		&userBook.Book.Author.ID,
		&userBook.Book.Genre.ID,
		&userBook.Book.Author.Name,
		&userBook.Book.Genre.Name,
		&userBook.User.ID,
		&userBook.User.Name,
		&userBook.User.ProfilePicture)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &userBook, nil
}

func FindUserBooksByUserID(ctx context.Context, userID int) ([]*UserBook, error) {
	// Find user books by user id and join with book and user
	rows, err := db.Conn.Query(ctx, `SELECT
		user_book.id,
		user_book.user_id,
		user_book.book_id,
		user_book.start_date,
		user_book.end_date,
		user_book.book_copy_id,
		book.id,
		book.title,
		book.subtitle,
		book.description,
		book.published_year,
		book.author_id,
		book.genre_id,
		author.id,
		author.name,
		genre.id,
		genre.name,
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.user_book
		JOIN public.book ON user_book.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.genre ON book.genre_id = genre.id
		JOIN public.user ON user_book.user_id = user.id
	WHERE
		user_book.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userBooks []*UserBook
	for rows.Next() {
		var userBook UserBook
		userBook.Book = &Book{}
		userBook.Book.Author = &Author{}
		userBook.Book.Genre = &Genre{}
		userBook.User = &User{}
		err = rows.Scan(
			&userBook.ID,
			&userBook.UserID,
			&userBook.BookID,
			&userBook.StartDate,
			&userBook.EndDate,
			&userBook.BookCopyID,
			&userBook.Book.ID,
			&userBook.Book.Title,
			&userBook.Book.Subtitle,
			&userBook.Book.Description,
			&userBook.Book.PublishedYear,
			&userBook.Book.Author.ID,
			&userBook.Book.Genre.ID,
			&userBook.Book.Author.Name,
			&userBook.Book.Genre.Name,
			&userBook.User.ID,
			&userBook.User.Name,
			&userBook.User.ProfilePicture)
		if err != nil {
			return nil, err
		}
		userBooks = append(userBooks, &userBook)
	}
	return userBooks, nil
}

func CreateUserBook(ctx context.Context, userBook *UserBook) (*UserBook, error) {
	// Create user book
	err := db.Conn.QueryRow(ctx, `INSERT INTO public.user_book (user_id, book_id, start_date, end_date) VALUES ($1, $2, $3, $4) RETURNING id`,
		userBook.UserID, userBook.BookID, userBook.StartDate, userBook.EndDate).Scan(&userBook.ID)
	if err != nil {
		return nil, err
	}
	return userBook, nil
}

func UpdateUserBook(ctx context.Context, userBook *UserBook) (*UserBook, error) {
	// Update user book
	_, err := db.Conn.Exec(ctx, `UPDATE public.user_book SET start_date = $1, end_date = $2 WHERE id = $3`,
		userBook.StartDate, userBook.EndDate, userBook.ID)
	if err != nil {
		return nil, err
	}
	return userBook, nil
}

func DeleteUserBook(ctx context.Context, id int) error {
	_, err := db.Conn.Exec(ctx, `DELETE FROM public.user_book WHERE id = $1`, id)
	return err
}
