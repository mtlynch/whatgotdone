// +build dev

package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s defaultServer) enableCsrf(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		h(w, r)
	}
}
