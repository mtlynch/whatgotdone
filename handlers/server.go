package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/auth"
	"github.com/mtlynch/whatgotdone/datastore"
)

type Server interface {
	Router() *mux.Router
}

func New() Server {
	s := defaultServer{
		authenticator: auth.New(),
		datastore:     datastore.New(),
		router:        mux.NewRouter(),
		csrfMiddleware: csrf.Protect(
			// TODO: Replace with an environment variable.
			[]byte("oh.38oV>iQy@E$Om15~>"),
			csrf.CookieName("csrf_base"),
			csrf.Path("/api/"),
			csrf.Secure(false)),
	}
	s.routes()
	return s
}

type httpMiddlewareHandler func(http.Handler) http.Handler

type defaultServer struct {
	authenticator  auth.Authenticator
	datastore      datastore.Datastore
	router         *mux.Router
	csrfMiddleware httpMiddlewareHandler
}

func (s defaultServer) Router() *mux.Router {
	return s.router
}
