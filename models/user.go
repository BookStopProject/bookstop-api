package models

import (
	"bookstop/db"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type User struct {
	ID             int       `json:"id"`
	OauthID        string    `json:"-"`
	Name           string    `json:"name"`
	Bio            *string   `json:"bio"`
	ProfilePicture *string   `json:"profilePicture"`
	CreationTime   time.Time `json:"creationTime"`
	Credit         int       `json:"credit"`
}

func FindUserByID(ctx context.Context, id int) (*User, error) {
	var user User
	err := db.Conn.QueryRow(ctx, `SELECT
		id,
		oauth_id,
		name,
		bio,
		profile_picture,
		creation_time,
		credit
	FROM
		public.user
	WHERE
		id = $1
`, id).Scan(
		&user.ID,
		&user.OauthID,
		&user.Name,
		&user.Bio,
		&user.ProfilePicture,
		&user.CreationTime,
		&user.Credit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func FindUserByOauthID(ctx context.Context, oauthID string) (*User, error) {
	var user User
	err := db.Conn.QueryRow(ctx, `SELECT 
		id,
		oauth_id,
		name,
		bio,
		profile_picture,
		creation_time,
		credit
	FROM
		public.user
	WHERE
		oauth_id = $1`, oauthID).Scan(
		&user.ID,
		&user.OauthID,
		&user.Name,
		&user.Bio,
		&user.ProfilePicture,
		&user.CreationTime,
		&user.Credit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func UpdateUser(ctx context.Context, user *User) (*User, error) {
	_, err := db.Conn.Exec(ctx, `UPDATE public.user SET
		name = $1,
		bio = $2
	WHERE
		id = $3`, user.Name, user.Bio, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(ctx context.Context, name string, oauthID string, email string, profilePicture *string) (*User, error) {
	var user User
	err := db.Conn.QueryRow(ctx, `INSERT INTO public.user (
		name,
		oauth_id,
		email,
		profile_picture
	) VALUES ($1, $2, $3, $4) RETURNING
		id,
		oauth_id,
		name,
		bio,
		profile_picture,
		creation_time,
		credit`, name, oauthID, email, profilePicture).Scan(
		&user.ID,
		&user.OauthID,
		&user.Name,
		&user.Bio,
		&user.ProfilePicture,
		&user.CreationTime,
		&user.Credit)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
