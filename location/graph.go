package location

import (
	"bookstop/graph/model"
	"strconv"

	"github.com/jackc/pgtype"
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
