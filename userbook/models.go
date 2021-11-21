package userbook

import (
	"bookstop/book"
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type UserBook struct {
	ID         pgtype.Int4
	UserID     pgtype.Int4
	BookID     pgtype.Varchar
	StartedAt  pgtype.Date
	EndedAt    pgtype.Date
	IDOriginal pgtype.Int4
}

const allSelects = "id, user_id, book_id, started_at, ended_at, id_original"

func verifyDates(startedAt *string, endedAt *string) error {
	var startedAtTime *time.Time
	var endedAtTime *time.Time

	if startedAt != nil {
		t, err := time.Parse("2006-01-02", *startedAt)
		if err != nil {
			return err
		}
		startedAtTime = &t
	}

	if endedAt != nil {
		t, err := time.Parse("2006-01-02", *endedAt)
		if err != nil {
			return err
		}
		endedAtTime = &t
	}

	if startedAtTime != nil && endedAtTime != nil {
		if (*startedAtTime).After(*endedAtTime) {
			return errors.New("start date cannot be after end date")
		}
	}

	return nil
}

func scanRow(row *pgx.Row) (*UserBook, error) {
	ub := &UserBook{}
	err := (*row).Scan(
		&ub.ID,
		&ub.UserID,
		&ub.BookID,
		&ub.StartedAt,
		&ub.EndedAt,
		&ub.IDOriginal,
	)
	if err != nil {
		return nil, err
	}
	return ub, nil
}

func scanRows(rows *pgx.Rows) (userBooks []*UserBook, errs []error) {
	for (*rows).Next() {
		ub := &UserBook{}
		err := (*rows).Scan(
			&ub.ID,
			&ub.UserID,
			&ub.BookID,
			&ub.StartedAt,
			&ub.EndedAt,
			&ub.IDOriginal,
		)
		if err != nil {
			errs = append(errs, err)
			userBooks = append(userBooks, nil)
		} else {
			errs = append(errs, nil)
			userBooks = append(userBooks, ub)
		}
	}
	return
}

func FindByID(ctx context.Context, id int) (*UserBook, error) {
	row := db.Conn.QueryRow(ctx, "SELECT "+allSelects+" FROM public.user_book WHERE id = $1", id)
	ub, err := scanRow(&row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return ub, nil
}

func Create(ctx context.Context, userID int, bookID string, startedAt *string, endedAt *string) (*UserBook, error) {
	// Verify book available
	b, err := book.FindByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("book not found")
	}
	// Verify date inputs
	err = verifyDates(startedAt, endedAt)
	if err != nil {
		return nil, err
	}

	row := db.Conn.QueryRow(ctx, `INSERT INTO public.user_book(
		user_id, book_id, started_at, ended_at)
		VALUES ($1, $2, $3, $4) RETURNING `+allSelects, userID, bookID, startedAt, endedAt)

	ub, err := scanRow(&row)

	if err != nil {
		return nil, err
	}
	return ub, nil
}

func IsOwner(ctx context.Context, userID int, id int) (bool, error) {
	var ubUserID int
	err := db.Conn.QueryRow(ctx, "SELECT user_id FROM public.user_book WHERE id = $1", id).Scan(&ubUserID)
	if err != nil {
		return false, err
	}
	return ubUserID == userID, nil
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*UserBook, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user_book WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result, errs := scanRows(&rows)

	idToPos := make(map[int]int)

	for i, ub := range result {
		idToPos[int(ub.ID.Int)] = i
	}

	sortedResult := make([]*UserBook, len(ids))
	sortedErrs := make([]error, len(ids))

	for i, id := range ids {
		pos := idToPos[id]
		sortedResult[i] = result[pos]
		sortedErrs[i] = errs[pos]
	}

	return sortedResult, sortedErrs
}

func FindManyByUserID(ctx context.Context, userID int) ([]*UserBook, []error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user_book WHERE user_id = $1 ORDER BY id DESC", userID)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	return scanRows(&rows)
}

func UpdateByID(ctx context.Context, id int, startedAt *string, endedAt *string) (*UserBook, error) {
	if startedAt == nil && endedAt == nil {
		return nil, errors.New("must provide at least one update value")
	}

	verifyDates(startedAt, endedAt)

	row := db.Conn.QueryRow(ctx, "UPDATE public.user_book SET started_at = $2, ended_at = $3 WHERE id = $1 RETURNING "+allSelects, id, startedAt, endedAt)

	return scanRow(&row)
}

func DeleteByID(ctx context.Context, id int) (bool, error) {
	rows, err := db.Conn.Query(ctx, "DELETE FROM public.user_book WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return true, nil
}
