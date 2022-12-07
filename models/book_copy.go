package models

import (
	"bookstop/db"
	"context"
)

type BookCondition string

const (
	BookConditionNew        BookCondition = "new"
	BookConditionLikeNew    BookCondition = "like_new"
	BookConditionGood       BookCondition = "good"
	BookConditionAcceptable BookCondition = "acceptable"
)

const (
	BookConditionNewMultiplier        = 1.0
	BookConditionLikeNewMultiplier    = 0.9
	BookConditionGoodMultiplier       = 0.7
	BookConditionAcceptableMultiplier = 0.5
)

type BookCopy struct {
	ID         int           `json:"id"`
	BookID     int           `json:"bookId"`
	Condition  BookCondition `json:"condition"`
	Book       *Book         `json:"book"`
	LocationID *int          `json:"locationId"`
	Location   *Location     `json:"location"`
}

func FindBookCopyOwners(ctx context.Context, bookCopyID int) ([]*User, error) {
	var users []*User
	rows, err := db.Conn.Query(ctx, `SELECT
		"user".id,
		"user".name,
		"user".profile_picture
	FROM
		public.user_book
		JOIN public."user" ON user_book.book_copy_id = "user".id
	WHERE
		user_book.book_copy_id = $1`, bookCopyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.ProfilePicture)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func FindBookCopiesByBookIdWithLocation(ctx context.Context, bookID int) ([]*BookCopy, error) {
	var bookCopies []*BookCopy
	// find book copy with location_id and join with location
	// and join with book and join with author
	rows, err := db.Conn.Query(ctx, `SELECT
		book_copy.id,
		book_copy.book_id,
		book_copy.condition,
		book_copy.location_id,
		book.id,
		book.title,
		book.subtitle,
		book.author_id,
		author.id,
		author.name,
		location.id,
		location.name,
		location.address
	FROM
		public.book_copy
		JOIN public.book ON book_copy.book_id = book.id
		JOIN public.author ON book.author_id = author.id
		JOIN public.location ON book_copy.location_id = location.id
	WHERE 
		book_copy.location_id IS NOT NULL
		AND book_copy.book_id = $1`, bookID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		bookCopy := &BookCopy{}
		bookCopy.Book = &Book{}
		bookCopy.Book.Author = &Author{}
		bookCopy.Location = &Location{}
		err := rows.Scan(
			&bookCopy.ID,
			&bookCopy.BookID,
			&bookCopy.Condition,
			&bookCopy.LocationID,
			&bookCopy.Book.ID,
			&bookCopy.Book.Title,
			&bookCopy.Book.Subtitle,
			&bookCopy.Book.AuthorID,
			&bookCopy.Book.Author.ID,
			&bookCopy.Book.Author.Name,
			&bookCopy.Location.ID,
			&bookCopy.Location.Name,
			&bookCopy.Location.Address,
		)
		if err != nil {
			return nil, err
		}
		bookCopies = append(bookCopies, bookCopy)
	}

	return bookCopies, nil
}
