// +build !dev

package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func getCsrfSeed() string {
	csrfSeed := os.Getenv("CSRF_SECRET_SEED")
	if csrfSeed == "" {
		log.Fatalf("CSRF_SECRET_SEED environment variable must be set")
	}
	return csrfSeed
}

func (s defaultServer) enableCsrf(h http.Handler) http.Handler {
	log.Printf("CSRF prod mode")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		log.Printf("r.RemoteAddr: %v", r.RemoteAddr)
		log.Printf("r.URL.Scheme: %v", r.URL.Scheme)
		log.Printf("r.URL: %+v", r.URL)
		log.Printf("r.Host: %v", r.Host)
		for key, values := range r.Header {
			for _, value := range values {
				log.Printf("%s: %v", key, value)
			}
		}
		s.csrfMiddleware(h).ServeHTTP(w, r)
	})
}
