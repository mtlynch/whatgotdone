package handlers

import (
	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/datastore"
)

type Server interface {
	Router() *mux.Router
}

func New() Server {
	s := defaultServer{
		datastore: datastore.New(),
		router:    mux.NewRouter(),
	}
	s.routes()
	return s
}

type defaultServer struct {
	datastore datastore.Datastore
	router    *mux.Router
}

func (s defaultServer) Router() *mux.Router {
	return s.router
}
