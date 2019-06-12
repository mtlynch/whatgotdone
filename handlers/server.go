package handlers

import (
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
	}
	s.routes()
	return s
}

type defaultServer struct {
	authenticator auth.Authenticator
	datastore     datastore.Datastore
	router        *mux.Router
}

func (s defaultServer) Router() *mux.Router {
	return s.router
}
