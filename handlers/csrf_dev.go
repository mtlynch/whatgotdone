// +build dev

package handlers

import (
	"net/http"
)

func (s defaultServer) enableCsrf(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
