// +build staging

package handlers

import (
	"net/http"
)

func (s defaultServer) enableCsp(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The staging environment needs the standard CSP policy + unsafe-eval
		// because it uses a frontend build with the "dev" tag which causes Vue to
		// generate JS that requires unsafe-eval.
		cspStaging := contentSecurityPolicy() + "; unsafe-eval"
		w.Header().Set("Content-Security-Policy", cspStaging)

		h(w, r)
	}
}
