package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"bookstop/db"
)

type TradeIn struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userId"`
	BookCopyID   int       `json:"bookCopyId"`
	Credit       int       `json:"credit"`
	CreationTime time.Time `json:"creationTime"`
	Book         *Book     `json:"book"`
}

func FindTradeInsByUserID(ctx context.Context, userId int) ([]*TradeIn, error) {
	// Find trade in by user id and join with book
	// and join with author
	rows, err := db.Conn.Query(ctx, `SELECT
		ti.id,
		ti.user_id,
		ti.book_copy_id,
		ti.credit,
		ti.creation_time,
		b.id,
		b.title,
		b.subtitle,
		b.published_year,
		a.id,
		a.name
	FROM trade_in ti
		JOIN book_copy bc ON ti.book_copy_id = bc.id
		JOIN book b ON bc.book_id = b.id
		JOIN author a ON b.author_id = a.id
	WHERE ti.user_id = $1
	`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tradeIns []*TradeIn

	for rows.Next() {
		var tradeIn TradeIn
		tradeIn.Book = &Book{}
		tradeIn.Book.Author = &Author{}

		err := rows.Scan(
			&tradeIn.ID,
			&tradeIn.UserID,
			&tradeIn.BookCopyID,
			&tradeIn.Credit,
			&tradeIn.CreationTime,
			&tradeIn.Book.ID,
			&tradeIn.Book.Title,
			&tradeIn.Book.Subtitle,
			&tradeIn.Book.PublishedYear,
			&tradeIn.Book.Author.ID,
			&tradeIn.Book.Author.Name,
		)

		if err != nil {
			return nil, err
		}

		tradeIns = append(tradeIns, &tradeIn)
	}

	return tradeIns, nil
}

func DoTradeIn(ctx context.Context, userBookID int, condition BookCondition, locationID int) (*TradeIn, error) {
	// TODO: implement procedure
	// This procedure should:
	// 1) Create a book copy if user book does not have one and link it to the user book
	// 2) Update the book copy condition and location
	// 3) Create a trade in for that book copy. The credit will be equal to
	//		the book trade in value * condition multiplier (see book_copy.go).
	// 4) Add the credit to the user's credit balance
	// 5) Return the trade in

	var tradeIn *TradeIn
	tradeIn.Book = new(Book)
	var err error

	getBookCopyIDQuery := `SELECT book_copy_id, user_id FROM user_book WHERE id=$1`

	getUserBookQuery := `SELECT * FROM book JOIN user_book ub ON book.id = ub.book_id`

	updateBookCopyQuery := `UPDATE book_copy SET condition=$1 && location_id=$2 WHERE id=$3`

	insertTradeInQuery := `INSERT INTO trade_in(id, user_id, book_copy_id, credit) VALUES ($1, $2, $3, $4)`

	if err := db.Conn.QueryRow(ctx, getUserBookQuery).Scan(
		&tradeIn.Book.ID,
		&tradeIn.Book.Title,
		&tradeIn.Book.Subtitle,
		&tradeIn.Book.PublishedYear,
		&tradeIn.Book.Author.ID,
		&tradeIn.Book.Author.Name,
	); err != nil {
		return nil, err
	}

	err = db.Conn.QueryRow(ctx, getBookCopyIDQuery).Scan(&tradeIn.BookCopyID, &tradeIn.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		insertQuery := `INSERT INTO book_copy(id, book_id, condition, location_id) VALUES ($1, $2, $3, $4)`
		if _, err := db.Conn.Exec(ctx, insertQuery, tradeIn.BookCopyID, tradeIn.Book.ID, condition, locationID);  err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}

	_, err = db.Conn.Exec(ctx, updateBookCopyQuery, condition, locationID, tradeIn.BookCopyID)
	if err != nil {
		return nil, err
	}

	_, err = db.Conn.Exec(ctx, insertTradeInQuery, tradeIn.BookCopyID, tradeIn.UserID, tradeIn.BookCopyID, float64(tradeIn.Book.TradeinCredit) * getBookConditionValue(condition))
	if err != nil {
		return nil, err
	}


	return tradeIn, nil
}

func getBookConditionValue(bc BookCondition) float64 {
	switch bc {
	case BookConditionNew:
		return BookConditionNewMultiplier
	case BookConditionLikeNew:
		return BookConditionLikeNewMultiplier
	case BookConditionGood:
		return BookConditionGoodMultiplier
	case BookConditionAcceptable:
		return BookConditionAcceptableMultiplier
	}
}
