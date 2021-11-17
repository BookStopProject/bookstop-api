package userbook

import (
	"bookstop/book"
	"bookstop/db"
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type UserBook struct {
	ID         pgtype.Int4
	UserId     pgtype.Int4
	BookId     pgtype.Varchar
	StartedAt  pgtype.Date
	EndedAt    pgtype.Date
	IDOriginal pgtype.Int4
}

const allSelects = "id, user_id, book_id, started_at, ended_at, id_original"

func FindById(ctx context.Context, id int) (*UserBook, error) {
	ub := &UserBook{}

	err := db.Conn.QueryRow(ctx, "SELECT "+allSelects+" FROM public.user_book WHERE id = $1", id).Scan(
		&ub.ID,
		&ub.UserId,
		&ub.BookId,
		&ub.StartedAt,
		&ub.EndedAt,
		&ub.IDOriginal,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return ub, nil
}

func addRowsToResult(rows *pgx.Rows, userBooks *[]*UserBook, errs *[]error) {
	for (*rows).Next() {
		ub := &UserBook{}
		err := (*rows).Scan(
			&ub.ID,
			&ub.UserId,
			&ub.BookId,
			&ub.StartedAt,
			&ub.EndedAt,
			&ub.IDOriginal,
		)
		if err != nil {
			ub = nil
		}
		*errs = append(*errs, err)
		*userBooks = append(*userBooks, ub)
	}
}

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

func Create(ctx context.Context, userId int, bookId string, startedAt *string, endedAt *string) (*UserBook, error) {
	ub := &UserBook{}

	// Verify book available
	b, err := book.FindById(ctx, bookId)
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

	err = db.Conn.QueryRow(ctx, `INSERT INTO public.user_book(
		user_id, book_id, started_at, ended_at)
		VALUES ($1, $2, $3, $4) RETURNING `+allSelects, userId, bookId, startedAt, endedAt).Scan(
		&ub.ID,
		&ub.UserId,
		&ub.BookId,
		&ub.StartedAt,
		&ub.EndedAt,
		&ub.IDOriginal,
	)
	if err != nil {
		return nil, err
	}
	return ub, nil
}

func IsOwner(ctx context.Context, userId int, id int) (bool, error) {
	var ubUserId int
	err := db.Conn.QueryRow(ctx, "SELECT user_id FROM public.user_book WHERE id = $1", id).Scan(&ubUserId)
	if err != nil {
		return false, err
	}
	return ubUserId == userId, nil
}

func FindManyByIds(ctx context.Context, ids []int) (userBooks []*UserBook, errs []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user_book WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)

	if err != nil {
		log.Panicln(err)
	}

	defer rows.Close()
	addRowsToResult(&rows, &userBooks, &errs)

	return
}

func FindManyByUserId(ctx context.Context, userId int) (userBooks []*UserBook, errs []error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user_book WHERE user_id = $1", userId)

	if err != nil {
		log.Panicln(err)
	}

	defer rows.Close()
	addRowsToResult(&rows, &userBooks, &errs)

	return
}

func UpdateById(ctx context.Context, id int, startedAt *string, endedAt *string) (*UserBook, error) {
	if startedAt == nil && endedAt == nil {
		return nil, errors.New("must provide at least one update value")
	}

	verifyDates(startedAt, endedAt)
	ub := UserBook{}

	err := db.Conn.QueryRow(ctx, "UPDATE public.user_book SET started_at = $2, ended_at = $3 WHERE id = $1 RETURNING "+allSelects, id, startedAt, endedAt).Scan(
		&ub.ID,
		&ub.UserId,
		&ub.BookId,
		&ub.StartedAt,
		&ub.EndedAt,
		&ub.IDOriginal,
	)

	if err != nil {
		return nil, err
	}

	return &ub, nil
}

func DeleteById(ctx context.Context, id int) (bool, error) {
	_, err := db.Conn.Query(ctx, "DELETE FROM public.user_book WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	return true, nil
}
