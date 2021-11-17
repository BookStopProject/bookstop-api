package auth

import (
	"bookstop/db"
	"bookstop/user"
	"context"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	redisAuthKeyPrefix = "auth:"
	authCodeCtxKey     = "authCode"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCode := r.Header.Get("authorization")

		ctx := context.WithValue(r.Context(), authCodeCtxKey, authCode)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func ForContext(ctx context.Context) (*user.User, error) {
	authCode, _ := ctx.Value(authCodeCtxKey).(string)
	if authCode == "" {
		return nil, nil
	}

	userId, err := db.RDB.Get(ctx, redisAuthKeyPrefix+authCode).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}

	return user.FindById(ctx, userIdInt)
}

func signIn(ctx context.Context, profile *GoogleProfileResponse) (string, error) {
	prevAuthCode, _ := ctx.Value(authCodeCtxKey).(string)
	if prevAuthCode != "" {
		// sign out previous auth
		signOut(ctx, prevAuthCode)
	}

	userId, err := user.FindIdByOauthId(ctx, profile.Id)
	if err != nil {
		return "", err
	}

	if userId == nil {
		// Create new user
		u, err := user.Create(ctx, profile.Name, profile.Id, profile.Email, profile.Picture)
		if err != nil {
			return "", err
		}
		newUserId := int(u.ID.Int)
		userId = &newUserId
	}

	authCode, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	_, err = db.RDB.Set(ctx, redisAuthKeyPrefix+authCode, *userId, 0).Result()

	if err != nil {
		return "", err
	}

	return authCode, err
}

func signOut(ctx context.Context, authCode string) (int64, error) {
	return db.RDB.Del(ctx, redisAuthKeyPrefix+authCode).Result()
}
