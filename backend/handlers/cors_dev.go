//go:build dev

package handlers

import (
	"net/http"
)

// enableCors sets the CORS headers so that when we run the frontend in dev
// mode, it can still communicate with the backend server.
func (s defaultServer) enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = r.Host
		}
		if origin == "" {
			http.Error(w, "(dev mode) Request needs a Host or Origin header", http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		h.ServeHTTP(w, r)
	})
}
