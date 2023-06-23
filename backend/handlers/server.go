package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/auth"
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/gcs"
)

// Server handles HTTP requests for the What Got Done backend.
type Server interface {
	Router() *mux.Router
}

// New creates a new What Got Done server with all the state it needs to
// satisfy HTTP requests.
func New(store datastore.Datastore) Server {
	gcsClient, err := gcs.New()
	if err != nil {
		log.Printf("failed to load Google Cloud Storage client: %s", err)
		log.Printf("file upload functionality will be disabled")
		gcsClient = nil
	}
	s := defaultServer{
		authenticator:  auth.New(),
		datastore:      store,
		gcsClient:      gcsClient,
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
	gcsClient      *gcs.Client
	router         *mux.Router
	csrfMiddleware httpMiddlewareHandler
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
