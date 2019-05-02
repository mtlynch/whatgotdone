package handlers

import (
	"github.com/gorilla/mux"
	userkit "github.com/workpail/userkit-go"

	"github.com/mtlynch/whatgotdone/datastore"
)

type Server interface {
	Router() *mux.Router
}

func New() Server {
	ks := datastore.NewUserKitKeyStore()
	sk, err := ks.SecretKey()
	if err != nil {
		panic(err)
	}
	s := defaultServer{
		datastore:     datastore.New(),
		userKitClient: userkit.NewUserKit(sk),
		router:        mux.NewRouter(),
	}
	s.routes()
	return s
}

type defaultServer struct {
	datastore     datastore.Datastore
	userKitClient userkit.Client
	router        *mux.Router
}

func (s defaultServer) Router() *mux.Router {
	return s.router
}
