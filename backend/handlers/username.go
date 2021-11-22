package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/types"
)

const userKitAuthCookieName = "userkit_auth_token"

type contextKey struct {
	name string
}

var contextKeyUsername = &contextKey{"username"}

func (s defaultServer) requireAuthentication(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenCookie, err := r.Cookie(userKitAuthCookieName)
		if err != nil {
			log.Printf("failed to retrieve cookie from request: %v", err)
			http.Error(w, "No authentication cookie found", http.StatusForbidden)
			return
		}

		username, err := s.authenticator.UserFromAuthToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, "Authentication required", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUsername, username)
		fn(w, r.WithContext(ctx))
	}
}

func usernameFromContext(ctx context.Context) types.Username {
	u, ok := ctx.Value(contextKeyUsername).(types.Username)
	if !ok {
		panic("expected to find username in context, but found none")
	}
	return u
}
