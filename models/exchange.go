package models

import (
	"bookstop/db"
	"context"
)

func DoExchange(ctx context.Context, userID int, bookCopyIDs []int) (*Invoice, error) {

	_, err := db.Conn.Exec(ctx, `CALL do_exchange($1, $2)`,
		userID, bookCopyIDs)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
