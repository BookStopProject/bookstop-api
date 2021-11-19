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

const allSelects = "id, name, description, image_url, started_at, ended_at"

type Browse = model.Browse

func scanRow(row *pgx.Row) (*Browse, error) {
	br := &Browse{}
	var browseId int
	err := (*row).Scan(
		&browseId,
		&br.Name,
		&br.Description,
		&br.ImageURL,
		&br.StartedAt,
		&br.EndedAt,
	)
	br.ID = strconv.Itoa(browseId)
	if err != nil {
		return nil, err
	}
	return br, nil
}

func FindById(ctx context.Context, id int) (*Browse, error) {
	row := db.Conn.QueryRow(ctx, "SELECT "+allSelects+" FROM public.browse WHERE id = $1", id)

	br, err := scanRow(&row)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return br, nil
}

func FindAll(ctx context.Context, when *time.Time) (results []*Browse, err error) {
	var rows pgx.Rows
	query := "SELECT " + allSelects + " FROM public.browse"
	if when != nil {
		rows, err = db.Conn.Query(ctx, query+" WHERE $1 between started_at and ended_at", when)
	} else {
		rows, err = db.Conn.Query(ctx, query)
	}

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		br := &Browse{}
		var browseId int
		err = rows.Scan(
			&browseId,
			&br.Name,
			&br.Description,
			&br.ImageURL,
			&br.StartedAt,
			&br.EndedAt,
		)
		br.ID = strconv.Itoa(browseId)
		if err != nil {
			return
		}
		results = append(results, br)
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
