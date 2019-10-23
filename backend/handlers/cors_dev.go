// +build dev

package handlers

import (
	"net/http"
)

// enableCors sets the CORS headers so that when we run the frontend in dev
// mode (on localhost:8085), it can still communicate with the backend server.
func (s defaultServer) enableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8085")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Csrf-Token")
		h(w, r)
	}
}
