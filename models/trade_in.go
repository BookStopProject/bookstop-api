package models

import (
	"bookstop/db"
	"context"
	"time"
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
		b.subtitle
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
