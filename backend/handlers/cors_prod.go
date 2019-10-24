// +build !dev

package handlers

import "net/http"

// enableCors on production builds is a passthrough that does nothing.
func (s defaultServer) enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
