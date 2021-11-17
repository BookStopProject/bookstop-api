package user

import (
	"bookstop/db"
	"context"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type User struct {
	ID              pgtype.Int4
	CreatedAt       pgtype.Timestamp
	OauthId         pgtype.Varchar
	Email           pgtype.Varchar
	Name            pgtype.Varchar
	Description     pgtype.Varchar
	ProfileImageUrl pgtype.Varchar
}

const (
	maxLengthName        = 50
	maxLengthDescription = 160
)

const allSelects = "id, name, oauth_id, email, description, profile_image_url, created_at"

func Create(ctx context.Context, name string, oauthId string, email *string, picture *string) (*User, error) {
	user := &User{}
	err := db.Conn.QueryRow(ctx, "INSERT INTO public.user(name, oauth_id, email, profile_image_url) VALUES ($1, $2, $3, $4) RETURNING "+allSelects, name, oauthId, email, picture).Scan(&user.ID,
		&user.Name,
		&user.OauthId,
		&user.Email,
		&user.Description,
		&user.ProfileImageUrl,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindIdByOauthId(ctx context.Context, oauthId string) (*int, error) {
	var id int

	err := db.Conn.QueryRow(ctx, "SELECT id FROM public.user WHERE oauth_id = $1", oauthId).Scan(
		&id,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &id, err
}

func FindById(ctx context.Context, id int) (*User, error) {
	user := &User{}

	err := db.Conn.QueryRow(ctx, "SELECT "+allSelects+" FROM public.user WHERE id = $1", id).Scan(
		&user.ID,
		&user.Name,
		&user.OauthId,
		&user.Email,
		&user.Description,
		&user.ProfileImageUrl,
		&user.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func FindManyByIds(ctx context.Context, ids []int) (users []*User, errs []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)

	if err != nil {
		log.Panicln(err)
	}

	defer rows.Close()

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID,
			&user.Name,
			&user.OauthId,
			&user.Email,
			&user.Description,
			&user.ProfileImageUrl,
			&user.CreatedAt)
		if err != nil {
			user = nil
		}
		errs = append(errs, err)
		users = append(users, user)
	}

	return
}

func UpdateById(ctx context.Context, id int, name *string, description *string) (*User, error) {
	if name == nil && description == nil {
		return nil, errors.New("must provide at least one update value")
	}

	var fields []string
	var values []interface{}
	values = append(values, id)

	i := 1

	if name != nil {
		i += 1
		*name = strings.TrimSpace(*name)

		if len(*name) <= 0 || len(*name) > maxLengthName {
			return nil, errors.New("Name length must be between 1 and " + strconv.Itoa(maxLengthName))
		}

		fields = append(fields, "name = $"+strconv.Itoa(i))
		values = append(values, *name)
	}

	if description != nil {
		i += 1
		*description = strings.TrimSpace(*description)

		if len(*description) > maxLengthDescription {
			return nil, errors.New("Bio length must be less than " + strconv.Itoa(maxLengthDescription))
		}

		fields = append(fields, "description = $"+strconv.Itoa(i))
		values = append(values, *description)
	}

	user := User{}

	err := db.Conn.QueryRow(ctx, "UPDATE public.user SET "+strings.Join(fields, ",")+" WHERE id = $1 RETURNING "+allSelects, values...).Scan(
		&user.ID,
		&user.Name,
		&user.OauthId,
		&user.Email,
		&user.Description,
		&user.ProfileImageUrl,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
