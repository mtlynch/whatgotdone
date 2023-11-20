package handlers

import (
	"net/http"
	"time"
)

func (s defaultServer) logoutPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    userKitAuthCookieName,
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})
	}
}
