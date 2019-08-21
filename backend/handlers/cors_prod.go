// +build !dev

package handlers

import "net/http"

func (s defaultServer) enableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
}
