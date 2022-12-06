package models

import (
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Event struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	LocationID  int       `json:"locationId"`
	Location    *Location `json:"location"`
}

type EventInventoryEntry struct {
	ID      int `json:"id"`
	EventID int `json:"eventId"`
}

func FindEventByID(ctx context.Context, id int) (*Event, error) {
	var event Event
	event.Location = &Location{}
	// find event and join location
	err := db.Conn.QueryRow(ctx, `SELECT
			e.id,
			e.name,
			e.description,
			e.start_time,
			e.end_time,
			e.location_id,
			l.name,
			l.address
		FROM
			public.event e
		JOIN
			public.location l ON e.location_id = l.id
		WHERE
			e.id = $1
		`, id).Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.LocationID,
		&event.Location.Name,
		&event.Location.Address,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func FindAllEvents(ctx context.Context) ([]*Event, error) {
	rows, err := db.Conn.Query(ctx, `SELECT
			e.id,
			e.name,
			e.description,
			e.start_time,
			e.end_time,
			e.location_id,
			l.name,
			l.address
		FROM
			public.event e
		JOIN
			public.location l ON e.location_id = l.id
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []*Event
	for rows.Next() {
		var event Event
		event.Location = &Location{}
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.LocationID,
			&event.Location.Name,
			&event.Location.Address,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
