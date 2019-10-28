// +build !dev

package handlers

import "net/http"

func (s *defaultServer) hstsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Scheme == "https" || // If the url tells us https
			r.TLS != nil && r.TLS.HandshakeComplete { // if the request is using tls

			rw.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

			next.ServeHTTP(rw, r)
			return
		}

		r.URL.Scheme = "https"

		http.Redirect(rw, r, r.URL.String(), http.StatusMovedPermanently)
	})
}
