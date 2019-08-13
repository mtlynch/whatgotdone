// +build dev

package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s defaultServer) enableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://10.0.0.100:8081")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		h(w, r)
	}
}
