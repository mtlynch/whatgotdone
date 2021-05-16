package handlers

import (
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/types"
)

const userKitAuthCookieName = "userkit_auth_token"

func (s defaultServer) loggedInUser(r *http.Request) (types.Username, error) {
	tokenCookie, err := r.Cookie(userKitAuthCookieName)
	if err != nil {
		return "", err
	}
	return s.authenticator.UserFromAuthToken(tokenCookie.Value)
}
