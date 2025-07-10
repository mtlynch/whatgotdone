package handlers

import (
	"net/http"
)

func (s *defaultServer) forwardingAddressGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}

func (s *defaultServer) forwardingAddressPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}

func (s *defaultServer) forwardingAddressDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	}
}
