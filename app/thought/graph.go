package thought

import (
	"bookstop/graph/model"
	"strconv"

	"github.com/jackc/pgtype"
)

func ToGraph(t *Thought) *model.Thought {
	if t == nil {
		return nil
	}

	val := model.Thought{
		ID:        strconv.Itoa(int(t.ID.Int)),
		Text:      t.Text.String,
		UserID:    strconv.Itoa(int(t.UserID.Int)),
		CreatedAt: t.CreatedAt.Time,
	}

	if t.BookID.Status == pgtype.Present {
		val.BookID = &t.BookID.String
	}

	return &val
}
