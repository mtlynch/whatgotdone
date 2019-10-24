// +build !staging

package handlers

import (
	"net/http"
)

func (s defaultServer) enableCsp(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", contentSecurityPolicy())
		h(w, r)
	}
}
