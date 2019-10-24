// +build dev

package handlers

import (
	"net/http"
)

// enableCors sets the CORS headers so that when we run the frontend in dev
// mode, it can still communicate with the backend server.
func (s defaultServer) enableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == "" {
			http.Error(w, "Cross domain requests require Origin header", http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		h.ServeHTTP(w, r)
	})
}
