package thought

import (
	"bookstop/book"
	"bookstop/db"
	"context"
	"errors"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Thought struct {
	ID        pgtype.Int4
	UserID    pgtype.Int4
	Text      pgtype.Varchar
	CreatedAt pgtype.Timestamp
	BookID    pgtype.Varchar
}

const allSelects = "id, user_id, created_at, text, book_id"

func scanRow(row *pgx.Row) (*Thought, error) {
	t := &Thought{}
	err := (*row).Scan(
		&t.ID,
		&t.UserID,
		&t.CreatedAt,
		&t.Text,
		&t.BookID,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func scanRows(rows *pgx.Rows) ([]*Thought, error) {
	var thoughts []*Thought
	for (*rows).Next() {
		t := &Thought{}
		err := (*rows).Scan(
			&t.ID,
			&t.UserID,
			&t.CreatedAt,
			&t.Text,
			&t.BookID,
		)
		if err != nil {
			return nil, err
		}
		thoughts = append(thoughts, t)
	}
	return thoughts, nil
}

func IsOwner(ctx context.Context, userID int, id int) (bool, error) {
	var ubUserID int
	err := db.Conn.QueryRow(ctx, "SELECT user_id FROM public.thought WHERE id = $1", id).Scan(&ubUserID)
	if err != nil {
		return false, err
	}
	return ubUserID == userID, nil
}

func Create(ctx context.Context, userID int, text string, bookID *string) (*Thought, error) {
	if bookID != nil {
		b, _ := book.FindByID(ctx, *bookID)
		if b == nil {
			return nil, errors.New("cannot find book")
		}
	}
	row := db.Conn.QueryRow(ctx, "INSERT INTO public.thought(user_id, text, book_id) VALUES ($1, $2, $3) RETURNING "+allSelects, userID, text, bookID)
	return scanRow(&row)
}

func FindAll(ctx context.Context, limit int, before *int) ([]*Thought, error) {
	if before == nil {
		bf := 999999
		before = &bf
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.thought WHERE id < $1 ORDER BY id DESC LIMIT $2", before, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(&rows)
}

func FindManyByUserID(ctx context.Context, userID int, limit int, before *int) ([]*Thought, error) {
	if before == nil {
		bf := 999999
		before = &bf
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.thought WHERE user_id = $1 AND id < $2 ORDER BY id DESC LIMIT $3", userID, before, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(&rows)
}

func DeleteByID(ctx context.Context, id int) (bool, error) {
	rows, err := db.Conn.Query(ctx, "DELETE FROM public.thought WHERE id = $1", id)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return true, nil
}
