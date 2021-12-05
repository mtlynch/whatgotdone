//go:build dev || staging

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	manager manager
	router  *mux.Router
}

func NewServer(mgr manager) server {
	return server{
		manager: mgr,
		router:  mux.NewRouter(),
	}
}

func (s *server) resetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.manager.Reset()
		if err != nil {
			log.Printf("failed to reset datstore: %v", err)
			http.Error(w, fmt.Sprintf("Failed to reset datastore: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
