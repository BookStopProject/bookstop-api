package auth

import (
	"bookstop/db"
	"bookstop/models"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
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

func ForContext(ctx context.Context) (*models.User, error) {
	authCode, _ := ctx.Value(authCodeCtxKey).(string)
	if authCode == "" {
		return nil, nil
	}

	userID, err := db.RDB.Get(ctx, redisAuthKeyPrefix+authCode).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	userIDInt, _ := strconv.Atoi(userID)

	return models.FindUserByID(ctx, userIDInt)
}

func signIn(ctx context.Context, profile *GoogleProfileResponse) (string, error) {
	prevAuthCode, _ := ctx.Value(authCodeCtxKey).(string)
	if prevAuthCode != "" {
		// sign out previous auth
		signOut(ctx, prevAuthCode)
	}

	user, err := models.FindUserByOauthID(ctx, profile.ID)
	if err != nil {
		return "", err
	}

	if user == nil {
		// Create new user
		if profile.Email == nil {
			return "", errors.New("email is required")
		}
		user, err = models.CreateUser(ctx, profile.Name, profile.ID, *profile.Email, profile.Picture)
		if err != nil {
			return "", err
		}
	}

	authCode, err := gonanoid.New()
	if err != nil {
		return "", err
	}

	_, err = db.RDB.Set(ctx, redisAuthKeyPrefix+authCode, user.ID, 0).Result()

	if err != nil {
		return "", err
	}

	return authCode, err
}

func signOut(ctx context.Context, authCode string) (int64, error) {
	return db.RDB.Del(ctx, redisAuthKeyPrefix+authCode).Result()
}

func getHmacSecret() []byte {
	hmacSecret := os.Getenv("HMAC_SECRET")
	if hmacSecret == "" {
		log.Fatalln("No HMAC_SECRET env")
	}
	return []byte(hmacSecret)
}
