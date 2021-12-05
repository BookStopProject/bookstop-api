package user

import (
	"bookstop/db"
	"bookstop/graph/model"
	"context"
	"strconv"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func ToGraph(user *User) *model.User {
	if user == nil {
		return nil
	}

	val := model.User{
		ID:        strconv.Itoa(int(user.ID.Int)),
		Name:      user.Name.String,
		CreatedAt: user.CreatedAt.Time,
	}

	if user.Description.Status == pgtype.Present {
		val.Description = &user.Description.String
	}

	if user.ProfileImageUrl.Status == pgtype.Present {
		val.ProfileImageURL = &user.ProfileImageUrl.String
	}

	return &val
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*model.User, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, `SELECT `+queryFieldsAll+`
	FROM public.user
	WHERE id IN (`+db.ParamRefsStr(len(ids))+`)`, args...)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	resultMap := make(map[int]*model.User)

	for rows.Next() {
		u, err := scanRowAll(rows.(pgx.Row))
		if err != nil {
			panic(err)
		}
		resultMap[int(u.ID.Int)] = ToGraph(u)
	}

	result := make([]*model.User, len(ids))
	errors := make([]error, len(ids))

	for i, id := range ids {
		ub, ok := resultMap[id]
		if ok {
			result[i] = ub
		}
	}

	return result, errors
}
