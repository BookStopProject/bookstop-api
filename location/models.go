package location

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Location struct {
	ID          pgtype.Int4
	Name        pgtype.Varchar
	ParentName  pgtype.Varchar
	AddressLine pgtype.Varchar
}

const queryFieldsAll = "id, name, parent_name, address_line"

func scanRowAll(row pgx.Row) (locs *Location, errs error) {
	loc := &Location{}
	err := (row).Scan(
		&loc.ID, &loc.Name, &loc.ParentName, &loc.AddressLine,
	)

	if err != nil {
		return nil, err
	}

	return loc, nil
}

func FindAllLocations(ctx context.Context) ([]*model.Location, error) {
	rows, err := db.Conn.Query(ctx, `SELECT `+queryFieldsAll+`
	FROM public.location`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []*model.Location

	for rows.Next() {
		loc, err := scanRowAll(rows.(pgx.Row))

		if err != nil {
			return nil, err
		}

		results = append(results, ToGraph(loc))
	}

	return results, nil
}

func Create(ctx context.Context, name string, parentName string, addressLine string) (*Location, error) {
	row := db.Conn.QueryRow(ctx, `INSERT INTO public.location(
	name, parent_name, address_line)
	VALUES ($1, $2, $3)
	RETURNING `+queryFieldsAll, name, parentName, addressLine)

	loc, err := scanRowAll(row)

	if err != nil {
		return nil, err
	}
	return loc, nil
}
