package exchange

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"
)

const allSelects = "id, user_book_id_old, user_book_id_new, exchanged_at"

func FindExchangesByUserBookId(ctx context.Context, ubId int) ([]*model.Exchange, error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.exchange WHERE user_book_id_old = $1 OR user_book_id_new = $1 OR user_book_id_original = $1 ORDER BY id DESC", ubId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var exchanges []*model.Exchange

	for rows.Next() {
		exc := &model.Exchange{}
		var exchangeId int
		var userBookIdOld int
		var userBookIdNew int
		err = rows.Scan(
			&exchangeId,
			&userBookIdOld,
			&userBookIdNew,
			&exc.ExchangedAt,
		)

		exc.ID = strconv.Itoa(exchangeId)
		exc.UserBookIDOld = strconv.Itoa(userBookIdOld)
		exc.UserBookIDNew = strconv.Itoa(userBookIdNew)

		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, exc)
	}

	return exchanges, nil
}
