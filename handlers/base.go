package handlers

import "net/http"

func New() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/auth/google/login", oauthGoogleLogin)
    mux.HandleFunc("/auth/google/callback", oauthGoogleCallback)

    return mux
}
