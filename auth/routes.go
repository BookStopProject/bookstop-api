package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetConf() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("API_URL") + "/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

type GoogleProfileResponse struct {
	ID      string  `json:"id"`
	Email   *string `json:"email"`
	Name    string  `json:"name"`
	Picture *string `json:"picture"`
}

func callbackWithError(w http.ResponseWriter, r *http.Request, redirectUrl string, err error) {
	if err == nil {
		err = errors.New("unknown")
	}
	if redirectUrl == "" {
		redirectUrl = "/"
	}
	appUrl := os.Getenv("APP_URL")
	http.Redirect(w, r, appUrl+redirectUrl+"?auth_error="+err.Error(), http.StatusTemporaryRedirect)
}

func apiAuth(w http.ResponseWriter, r *http.Request) {
	url := GetConf().AuthCodeURL(r.URL.Query().Get("redirect_url"))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func apiCallback(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	errCode := r.URL.Query().Get("error")
	redirectUrl := r.URL.Query().Get("state")

	if errCode != "" {
		callbackWithError(w, r, redirectUrl, errors.New(errCode))
		return
	}

	ctx := r.Context()

	conf := GetConf()
	tok, err := conf.Exchange(ctx, authCode)
	if err != nil {
		callbackWithError(w, r, redirectUrl, err)
		return
	}
	client := conf.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		callbackWithError(w, r, redirectUrl, nil)
		return
	}
	defer resp.Body.Close()
	var profileResponse GoogleProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&profileResponse)

	if err != nil {
		callbackWithError(w, r, redirectUrl, nil)
		return
	}

	authToken, err := signIn(ctx, &profileResponse)

	if err != nil {
		callbackWithError(w, r, redirectUrl, err)
		return
	}

	appUrl := os.Getenv("APP_URL")
	http.Redirect(w, r, appUrl+redirectUrl+"?auth_token="+authToken, http.StatusTemporaryRedirect)
}

func Router(router *httprouter.Router) {
	(*router).HandlerFunc(http.MethodGet, "/auth", apiAuth)
	(*router).HandlerFunc(http.MethodGet, "/auth/callback", apiCallback)
}
