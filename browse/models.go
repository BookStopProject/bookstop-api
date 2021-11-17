package browse

import (
	"bookstop/book"
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
)

const allSelects = "id, name, description, image_url"

func FindById(ctx context.Context, id int) (*model.Browse, error) {
	browse := &model.Browse{}

	err := db.Conn.QueryRow(ctx, "SELECT "+allSelects+" FROM public.browse WHERE id = $1", id).Scan(
		browse.ID,
		browse.Name,
		browse.Description,
		browse.ImageURL,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return browse, nil
}

func FindAll(ctx context.Context, when *time.Time) (results []*model.Browse, err error) {
	query := "SELECT " + allSelects + " FROM public.browse"
	if when != nil {
		query += " WHERE $1 between started_at and ended_at"
	}
	rows, err := db.Conn.Query(ctx, query, when)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		browse := &model.Browse{}
		var browseId int
		err = rows.Scan(
			&browseId,
			&browse.Name,
			&browse.Description,
			&browse.ImageURL,
		)
		browse.ID = strconv.Itoa(browseId)
		if err != nil {
			return
		}
		results = append(results, browse)
	}

	return
}

func FindBooksByBrowseId(ctx context.Context, id int) (results []*book.Book, errs []error) {
	bookIds := []string{}

	rows, err := db.Conn.Query(ctx, "SELECT book_id FROM public.browse_book where browse_id = $1", id)

	if err != nil {
		errs = append(errs, err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var bookId string
		rows.Scan(&bookId)
		bookIds = append(bookIds, bookId)
	}

	results, errs = book.FindManyByIds(ctx, bookIds)
	return
}
