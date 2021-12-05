package userbook

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func ToGraph(userBook *UserBook) *model.UserBook {
	if userBook == nil {
		return nil
	}

	val := model.UserBook{
		ID:     strconv.Itoa(int(userBook.ID.Int)),
		BookID: userBook.BookID.String,
		UserID: strconv.Itoa(int(userBook.UserID.Int)),
	}

	if userBook.StartedAt.Status == pgtype.Present {
		startedAt := userBook.StartedAt.Time.Format("2006-01-02")
		val.StartedAt = &startedAt
	}

	if userBook.EndedAt.Status == pgtype.Present {
		endedAt := userBook.EndedAt.Time.Format("2006-01-02")
		val.EndedAt = &endedAt
	}

	if userBook.IDOriginal.Status == pgtype.Present {
		originalID := strconv.Itoa(int(userBook.IDOriginal.Int))
		val.OriginalUserBookID = &originalID
	}

	return &val
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*model.UserBook, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+queryFieldsAll+" FROM public.user_book WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	resultMap := make(map[int]*model.UserBook)

	for rows.Next() {
		ub, err := scanRow(rows.(pgx.Row))
		if err != nil {
			panic(err)
		}
		resultMap[int(ub.ID.Int)] = ToGraph(ub)
	}

	result := make([]*model.UserBook, len(ids))
	errors := make([]error, len(ids))

	for i, id := range ids {
		ub, ok := resultMap[id]
		if ok {
			result[i] = ub
		}
	}

	return result, errors
}
