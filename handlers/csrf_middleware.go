package handlers

import (
	"log"
	"os"

	"github.com/gorilla/csrf"
)

func newCsrfMiddleware() httpMiddlewareHandler {
	csrfSeed := os.Getenv("CSRF_SECRET_SEED")
	if csrfSeed == "" {
		log.Fatalf("CSRF_SECRET_SEED environment variable must be set")
	}
	return csrf.Protect(
		[]byte(csrfSeed),
		csrf.CookieName("csrf_base"),
		csrf.Path("/api/"),
		csrf.Secure(false))
}
