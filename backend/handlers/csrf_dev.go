//go:build dev

package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func getCsrfSeed() string {
	// In dev mode, use a hardcoded CSRF secret seed.
	return "dummy-dev-csrf-seed"
}

func (s defaultServer) enableCsrf(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In dev mode, set the X-CSRF-Token on responses, but don't require
		// requests to send a correct CSRF token back.
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		h.ServeHTTP(w, r)
	})
}
