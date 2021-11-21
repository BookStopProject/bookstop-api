package user

import (
	"bookstop/db"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type User struct {
	ID              pgtype.Int4
	CreatedAt       pgtype.Timestamp
	OauthID         pgtype.Varchar
	Email           pgtype.Varchar
	Name            pgtype.Varchar
	Description     pgtype.Varchar
	ProfileImageUrl pgtype.Varchar
	Credit          pgtype.Int4
}

const (
	maxLengthName        = 50
	maxLengthDescription = 160
)

const allSelects = "id, name, oauth_id, email, description, profile_image_url, created_at, credit"

func scanRow(row *pgx.Row) (*User, error) {
	user := &User{}
	err := (*row).Scan(
		&user.ID,
		&user.Name,
		&user.OauthID,
		&user.Email,
		&user.Description,
		&user.ProfileImageUrl,
		&user.CreatedAt,
		&user.Credit,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func scanRows(rows *pgx.Rows) (users []*User, errs []error) {
	for (*rows).Next() {
		user := &User{}
		err := (*rows).Scan(
			&user.ID,
			&user.Name,
			&user.OauthID,
			&user.Email,
			&user.Description,
			&user.ProfileImageUrl,
			&user.CreatedAt,
			&user.Credit,
		)
		if err != nil {
			errs = append(errs, err)
			users = append(users, nil)
		} else {
			errs = append(errs, nil)
			users = append(users, user)
		}
	}
	return
}

func Create(ctx context.Context, name string, oauthID string, email *string, picture *string) (*User, error) {
	row := db.Conn.QueryRow(ctx, "INSERT INTO public.user(name, oauth_id, email, profile_image_url) VALUES ($1, $2, $3, $4) RETURNING "+allSelects, name, oauthID, email, picture)
	return scanRow(&row)
}

func FindIDByOauthID(ctx context.Context, oauthID string) (*int, error) {
	var id int

	err := db.Conn.QueryRow(ctx, "SELECT id FROM public.user WHERE oauth_id = $1", oauthID).Scan(
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

func FindByID(ctx context.Context, id int) (*User, error) {
	row := db.Conn.QueryRow(ctx, "SELECT "+allSelects+" FROM public.user WHERE id = $1", id)

	user, err := scanRow(&row)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func FindAll(ctx context.Context) ([]*User, []error) {
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	return scanRows(&rows)
}

func LoadManyByIDs(ctx context.Context, ids []int) ([]*User, []error) {
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}
	rows, err := db.Conn.Query(ctx, "SELECT "+allSelects+" FROM public.user WHERE id IN ("+db.ParamRefsStr(len(ids))+")", args...)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	result, errs := scanRows(&rows)

	idToPos := make(map[int]int)

	for i, ub := range result {
		idToPos[int(ub.ID.Int)] = i
	}

	sortedResult := make([]*User, len(ids))
	sortedErrs := make([]error, len(ids))

	for i, id := range ids {
		pos := idToPos[id]
		sortedResult[i] = result[pos]
		sortedErrs[i] = errs[pos]
	}

	return sortedResult, sortedErrs
}

func UpdateByID(ctx context.Context, id int, name *string, description *string) (*User, error) {
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

	row := db.Conn.QueryRow(ctx, "UPDATE public.user SET "+strings.Join(fields, ",")+" WHERE id = $1 RETURNING "+allSelects, values...)

	return scanRow(&row)
}
