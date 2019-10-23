package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/auth"
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/firestore"
)

// Server handles HTTP requests for the What Got Done backend.
type Server interface {
	Router() *mux.Router
}

// New creates a new What Got Done server with all the state it needs to
// satisfy HTTP requests.
func New() Server {
	s := defaultServer{
		authenticator:  auth.New(),
		datastore:      firestore.New(),
		router:         mux.NewRouter(),
		csrfMiddleware: newCsrfMiddleware(),
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

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
