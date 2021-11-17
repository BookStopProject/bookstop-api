package user

import (
	"bookstop/graph/model"
	"strconv"

	"github.com/jackc/pgtype"
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
