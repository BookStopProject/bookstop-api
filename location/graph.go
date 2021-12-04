package location

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func ToGraph(location *Location) *model.Location {
	if location == nil {
		return nil
	}

	val := model.Location{
		ID:          strconv.Itoa(int(location.ID.Int)),
		Name:        location.Name.String,
		AddressLine: location.AddressLine.String,
	}

	if location.ParentName.Status == pgtype.Present {
		val.ParentName = &location.ParentName.String
	}

	return &val
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*model.Location, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+queryFieldsAll+" FROM public.location WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	resultMap := make(map[int]*model.Location)

	for rows.Next() {
		loc, err := scanRowAll(rows.(pgx.Row))
		if err != nil {
			panic(err)
		}
		resultMap[int(loc.ID.Int)] = ToGraph(loc)
	}

	result := make([]*model.Location, len(ids))
	errors := make([]error, len(ids))

	for i, id := range ids {
		ub, ok := resultMap[id]
		if ok {
			result[i] = ub
		}
	}

	return result, errors
}
