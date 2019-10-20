package handlers

import (
	"net/http"
)

const userKitAuthCookieName = "userkit_auth_token"

func (s defaultServer) loggedInUser(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(userKitAuthCookieName)
	if err != nil {
		return "", err
	}
	return s.authenticator.UserFromAuthToken(tokenCookie.Value)
}
