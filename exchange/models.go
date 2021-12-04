package exchange

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"
)

const queryFieldsAll = "id, user_book_id_old, user_book_id_new, exchanged_at"

func FindByUserBookID(ctx context.Context, ubID int) ([]*model.Exchange, error) {
	rows, err := db.Conn.Query(ctx, `SELECT `+queryFieldsAll+`
	FROM public.exchange
	WHERE WHERE user_book_id_old = $1 OR user_book_id_new = $1 OR user_book_id_original = $1
	ORDER BY id DESC`, ubID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var exchanges []*model.Exchange

	for rows.Next() {
		exc := &model.Exchange{}
		var exchangeID int
		var userBookIDOld int
		var userBookIDNew int
		err = rows.Scan(
			&exchangeID,
			&userBookIDOld,
			&userBookIDNew,
			&exc.ExchangedAt,
		)

		exc.ID = strconv.Itoa(exchangeID)
		exc.UserBookIDOld = strconv.Itoa(userBookIDOld)
		exc.UserBookIDNew = strconv.Itoa(userBookIDNew)

		if err != nil {
			return nil, err
		}

		exchanges = append(exchanges, exc)
	}

	return exchanges, nil
}
