package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const state_cookie = "oauthstate"
const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

var googleOauthConfig = &oauth2.Config{
    RedirectURL: "http://localhost:8000/auth/google/callback",
    ClientID: os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
    Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"},
    Endpoint: google.Endpoint,
}

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
    // Create oauthState cookie and redirect to Google authentication page
    oauthState := generateStateOauthCookie(w)
    authCodeURL := googleOauthConfig.AuthCodeURL(oauthState)
    http.Redirect(w, r, authCodeURL, http.StatusTemporaryRedirect)
}

// State is a token to protect the user from CSRF attacks. You must
// always provide a non-empty string and validate that it matches the
// the state query parameter on your redirect callback.
func generateStateOauthCookie(w http.ResponseWriter) string {
    expiration := time.Now().Add(365 * 24 * time.Hour)
    b := make([]byte, 16)
    rand.Read(b)
    state := base64.URLEncoding.EncodeToString(b)
    cookie := http.Cookie{Name: state_cookie, Value: state, Expires: expiration}
    http.SetCookie(w, &cookie)

    return state
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
    // Read oauthState from Cookie
    oauthState, _ := r.Cookie(state_cookie)

    if r.FormValue("state") != oauthState.Value {
        log.Println("Invalid oauth google state")
        http.Error(w, "Invalid oauth google state", 401)
        return
    }

    token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
    if err != nil {
        log.Println(err.Error())
        http.Error(w, err.Error(), 401)
        return
    }

    // TODO: get or create user in db
    // return response with token
    if r.Header.Get("Content-Type") == "application/json" {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(token)
        return
    }

    fmt.Fprintf(w, token.AccessToken)
}

