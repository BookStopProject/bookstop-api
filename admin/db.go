package admin

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5"
)

func getDbConn(r *http.Request) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, errors.New("invalid auth")
	}
	dbURLRaw := os.Getenv("DATABASE_URL")
	dbURL, err := url.Parse(dbURLRaw)
	if err != nil {
		return nil, err
	}
	dbURL.User = url.UserPassword(username, password)
	dbURLRaw = dbURL.String()
	conn, err = pgx.Connect(r.Context(), dbURLRaw)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
