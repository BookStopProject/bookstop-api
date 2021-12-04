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

const allSelects = "id, name, parent_name, address_line"

func scanRows(rows *pgx.Rows) (locs []*Location, errs []error) {
	for (*rows).Next() {
		loc := &Location{}
		err := (*rows).Scan(
			&loc.ID, &loc.Name, &loc.ParentName, &loc.AddressLine,
		)
		if err != nil {
			errs = append(errs, err)
			locs = append(locs, nil)
		} else {
			errs = append(errs, nil)
			locs = append(locs, loc)
		}
	}
	return
}

func FindAllLocations(ctx context.Context) ([]*model.Location, error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.location")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []*model.Location

	for rows.Next() {
		loc := &Location{}
		err = rows.Scan(
			&loc.ID, &loc.Name, &loc.ParentName, &loc.AddressLine)

		if err != nil {
			return nil, err
		}

		results = append(results, ToGraph(loc))
	}

	return results, nil
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*Location, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.location WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result, errs := scanRows(&rows)

	idToPos := make(map[int]int)

	for i, ub := range result {
		idToPos[int(ub.ID.Int)] = i
	}

	sortedResult := make([]*Location, len(ids))
	sortedErrs := make([]error, len(ids))

	for i, id := range ids {
		pos, ok := idToPos[id]
		if ok {
			sortedResult[i] = result[pos]
			sortedErrs[i] = errs[pos]
		}
	}

	return sortedResult, sortedErrs
}

func Create(ctx context.Context, name string, parentName string, addressLine string) (*Location, error) {
	loc := &Location{}
	err := db.Conn.QueryRow(ctx, `INSERT INTO public.location(name, parent_name, address_line) VALUES ($1, $2, $3) RETURNING `+allSelects, name, parentName, addressLine).Scan(
		&loc.ID, &loc.Name, &loc.ParentName, &loc.AddressLine,
	)
	if err != nil {
		return nil, err
	}
	return loc, nil
}
