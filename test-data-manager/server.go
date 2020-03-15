// +build dev integration

package main

import (
	"encoding/json"
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
		type response struct {
			Ok bool `json:"ok"`
		}
		resp := response{
			Ok: true,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
