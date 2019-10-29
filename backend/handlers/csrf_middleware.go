package handlers

import (
	"github.com/gorilla/csrf"
)

func newCsrfMiddleware() httpMiddlewareHandler {
	csrfSeed := getCsrfSeed()
	return csrf.Protect(
		[]byte(csrfSeed),
		// The _v suffix is just to prevent clients from re-using a previous version of the cookie.
		// When rev'ing the version, be sure to globally replace this name in the entire codebase.
		csrf.CookieName("csrf_base_v3"),
		csrf.Path("/"),
		csrf.Secure(false))
}
