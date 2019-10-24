// +build !dev

package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s defaultServer) enableCsrf(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))

		s.csrfMiddleware(h).ServeHTTP(w, r)
	})
}
