// +build dev

package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

func getCsrfSeed() string {
	// In dev mode, use a hardcoded CSRF secret seed.
	return "dummy-dev-csrf-seed"
}

func (s defaultServer) enableCsrf(h http.Handler) http.Handler {
	log.Printf("enabling CSRF protection in dev mode")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// In dev mode, set the X-CSRF-Token on responses, but don't require
		// requests to send a correct CSRF token back.
		w.Header().Set("X-CSRF-Token", csrf.Token(r))

		for key, values := range r.Header {
			for _, value := range values {
				log.Printf("%s: %v", key, value)
			}
		}
		h.ServeHTTP(w, r)
	})
}
