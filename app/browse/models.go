package browse

import (
	"bookstop/app/book"
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
)

type Browse = model.Browse

const queryFieldsAll = "id, name, description, image_url, started_at, ended_at"

func scanRowAll(row pgx.Row) (*Browse, error) {
	br := &Browse{}
	var browseID int
	err := row.Scan(
		&browseID,
		&br.Name,
		&br.Description,
		&br.ImageURL,
		&br.StartedAt,
		&br.EndedAt,
	)
	br.ID = strconv.Itoa(browseID)
	if err != nil {
		return nil, err
	}
	return br, nil
}

func FindByID(ctx context.Context, id int) (*Browse, error) {
	row := db.Conn.QueryRow(ctx, `SELECT `+queryFieldsAll+`
	FROM public.browse
	WHERE id = $1`, id)

	br, err := scanRowAll(row)

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
	query := "SELECT " + queryFieldsAll + " FROM public.browse"
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
		br, err := scanRowAll(rows.(pgx.Row))
		if err != nil {
			return nil, err
		}
		results = append(results, br)
	}

	return
}

func FindBooksByBrowseID(ctx context.Context, id int) (results []*book.Book, errs []error) {
	rows, err := db.Conn.Query(ctx, `SELECT book_id
	FROM public.browse_book
	WHERE browse_id = $1`, id)

	if err != nil {
		errs = append(errs, err)
		return
	}

	defer rows.Close()

	bookIDs := []string{}

	for rows.Next() {
		var bookID string
		rows.Scan(&bookID)
		bookIDs = append(bookIDs, bookID)
	}

	results, errs = book.LoadManyByIDs(ctx, bookIDs)
	return
}

func Create(ctx context.Context, name string, description string, startedAt string, endedAt string) (*Browse, error) {
	row := db.Conn.QueryRow(ctx, `INSERT INTO public.browse(
	name, description, started_at, ended_at)
	VALUES ($1, $2, $3, $4)
	RETURNING `+queryFieldsAll, name, description, startedAt, endedAt)
	return scanRowAll(row)
}

func UpdateByID(ctx context.Context, id int, name string, description string, startedAt string, endedAt string) (*Browse, error) {
	row := db.Conn.QueryRow(ctx, `UPDATE public.browse
	SET name = $2, description = $3, started_at = $4, ended_at = $5
	WHERE id = $1
	RETURNING `+queryFieldsAll, id, name, description, startedAt, endedAt)

	return scanRowAll(row)
}

func DeleteByID(ctx context.Context, id int) (bool, error) {
	_, err := db.Conn.Query(ctx, `DELETE FROM public.browse
	WHERE id= $1`, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddBooksByIDs(ctx context.Context, id int, bookIDs []string) (bool, error) {
	if len(bookIDs) <= 0 {
		return false, nil
	}

	args := make([]interface{}, len(bookIDs)+1)
	args[0] = id

	query := "INSERT INTO public.browse_book(book_id, browse_id) VALUES"

	_, errs := book.LoadManyByIDs(ctx, bookIDs)
	for i, bookID := range bookIDs {
		if errs[i] != nil {
			return false, errors.New("book " + bookIDs[i] + ": " + errs[i].Error())
		}
		query += " ($" + (strconv.Itoa(i+2) + ", $1)")

		bookID = strings.TrimSpace(bookID)
		if bookID == "" {
			return false, errors.New("empty book id")
		}
		args[i+1] = bookID

		if i < len(errs)-1 {
			query += ","
		}
	}

	rows, err := db.Conn.Query(ctx, query, args...)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	return true, nil
}

func DeleteBooksByIDs(ctx context.Context, id int, bookIDs []string) (bool, error) {
	if len(bookIDs) <= 0 {
		return false, nil
	}

	args := make([]interface{}, len(bookIDs)+1)

	for i, bookID := range bookIDs {
		bookID = strings.TrimSpace(bookID)
		if bookID == "" {
			return false, errors.New("empty book id")
		}
		args[i] = bookID
	}

	args[len(args)-1] = id

	query := `DELETE FROM public.browse_book
	WHERE book_id IN (` + db.ParamRefsStr(len(bookIDs)) + `)
	AND browse_id = $` + strconv.Itoa(len(bookIDs)+1)

	rows, err := db.Conn.Query(ctx, query, args...)

	if err != nil {
		return false, err
	}

	defer rows.Close()
	return true, nil
}
