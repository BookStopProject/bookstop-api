package event

import (
	"bookstop/db"
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Event struct {
	ID          pgtype.Int4
	Title       pgtype.Varchar
	Description pgtype.Varchar
	Href        pgtype.Varchar
	UserID      pgtype.Int4
	StartedAt   pgtype.Timestamp
	EndedAt     pgtype.Timestamp
}

const queryFieldsAll = "id, title, description, href, user_id, started_at, ended_at"

func scanRowAll(row pgx.Row) (*Event, error) {
	evt := &Event{}
	err := row.Scan(
		&evt.ID,
		&evt.Title,
		&evt.Description,
		&evt.Href,
		&evt.UserID,
		&evt.StartedAt,
		&evt.EndedAt,
	)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func FindAll(ctx context.Context) ([]*Event, error) {
	rows, err := db.Conn.Query(ctx, `SELECT `+queryFieldsAll+`
	FROM public.event`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []*Event

	for rows.Next() {
		evt, err := scanRowAll(rows.(pgx.Row))
		if err != nil {
			return nil, err
		}
		events = append(events, evt)
	}

	return events, nil
}

func Create(ctx context.Context, title string, description string, href string, userId string, startedAt string, endedAt string) (*Event, error) {
	row := db.Conn.QueryRow(ctx, `INSERT INTO public.event(
	title, description, href, user_id, started_at, ended_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING `+queryFieldsAll, title, description, href, userId, startedAt, endedAt)
	evt, err := scanRowAll(row)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func Update(ctx context.Context, id int, title string, description string, href string, userId int, startedAt time.Time, endedAt time.Time) (*Event, error) {
	evt := &Event{}
	row := db.Conn.QueryRow(ctx, `UPDATE public.event
	SET title=$2, description=$3, href=$4, user_id=$5, started_at=$6, ended_at=$7
	WHERE id = $1
	RETURNING `+queryFieldsAll, id, title, description, href, userId, startedAt, endedAt)
	evt, err := scanRowAll(row)
	if err != nil {
		return nil, err
	}
	return evt, nil
}

func Remove(ctx context.Context, id int) error {
	rows, err := db.Conn.Query(ctx, `DELETE FROM public.event
	WHERE id=$1`, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}
