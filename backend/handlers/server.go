package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/auth"
	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/gcs"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
)

// Server handles HTTP requests for the What Got Done backend.
type Server interface {
	Router() *mux.Router
}

// New creates a new What Got Done server with all the state it needs to
// satisfy HTTP requests.
func New() Server {
	var fetcher ga.MetricFetcher
	f, err := ga.New()
	if err != nil {
		log.Printf("Failed to load Google Analytics metrics fetcher: %s", err)
	} else {
		fetcher = f
	}

	gcsClient, err := gcs.New()
	if err != nil {
		log.Printf("Failed to load Google Cloud Storage client: %s", err)
		log.Printf("File upload functionality will be disabled")
		gcsClient = nil
	}
	ds := newDatastore()
	s := defaultServer{
		authenticator:          auth.New(),
		datastore:              ds,
		gcsClient:              gcsClient,
		router:                 mux.NewRouter(),
		csrfMiddleware:         newCsrfMiddleware(),
		googleAnalyticsFetcher: fetcher,
	}
	s.routes()
	return s
}

type httpMiddlewareHandler func(http.Handler) http.Handler

type defaultServer struct {
	authenticator          auth.Authenticator
	datastore              datastore.Datastore
	gcsClient              *gcs.Client
	router                 *mux.Router
	csrfMiddleware         httpMiddlewareHandler
	googleAnalyticsFetcher ga.MetricFetcher
}

// Router returns the underlying router interface for the server.
func (s defaultServer) Router() *mux.Router {
	return s.router
}
