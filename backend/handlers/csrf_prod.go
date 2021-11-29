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
	log.Printf("enabling CSRF protection in production mode")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		s.csrfMiddleware(h).ServeHTTP(w, r)
	})
}
