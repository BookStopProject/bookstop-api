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
	BookConditionNewScale        = 1.0
	BookConditionLikeNewScale    = 0.9
	BookConditionGoodScale       = 0.7
	BookConditionAcceptableScale = 0.5
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
