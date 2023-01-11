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

func (s defaultServer) populateAuthenticationContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(userKitAuthCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		username, err := s.authenticator.UserFromAuthToken(tokenCookie.Value)
		if err != nil {
			log.Printf("failed to get username from auth token: %v", err)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyUsername, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s defaultServer) requireAuthenticationForApi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := usernameFromContext(r.Context()); !ok {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s defaultServer) requireAuthenticationForView(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := usernameFromContext(r.Context()); !ok {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func usernameFromContext(ctx context.Context) (types.Username, bool) {
	u, ok := ctx.Value(contextKeyUsername).(types.Username)
	return u, ok
}

func mustGetUsernameFromContext(ctx context.Context) types.Username {
	u, ok := ctx.Value(contextKeyUsername).(types.Username)
	if !ok {
		panic("username not found in context")
	}
	return u
}
