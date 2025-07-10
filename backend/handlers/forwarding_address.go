package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/handlers/parse"
)

type forwardingAddressResponse struct {
	ForwardingUrl string `json:"forwardingUrl"`
}

type forwardingAddressRequest struct {
	ForwardingUrl string `json:"forwardingUrl"`
}

func (s *defaultServer) forwardingAddressGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mustGetUsernameFromContext(r.Context())

		address, err := s.datastore.GetForwardingAddress(username)
		if err != nil {
			if _, ok := err.(datastore.ForwardingAddressNotFoundError); ok {
				// Return empty forwarding URL if none is set
				respondOK(w, forwardingAddressResponse{ForwardingUrl: ""})
				return
			}
			log.Printf("failed to retrieve forwarding address for user %s: %v", username, err)
			http.Error(w, "Failed to retrieve forwarding address", http.StatusInternalServerError)
			return
		}

		respondOK(w, forwardingAddressResponse{ForwardingUrl: string(address)})
	}
}

func (s *defaultServer) forwardingAddressPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mustGetUsernameFromContext(r.Context())

		var req forwardingAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("failed to decode forwarding address request: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.ForwardingUrl == "" {
			http.Error(w, "Forwarding URL cannot be empty", http.StatusBadRequest)
			return
		}

		address, err := parse.ForwardingAddress(req.ForwardingUrl)
		if err != nil {
			log.Printf("failed to parse forwarding address %s: %v", req.ForwardingUrl, err)
			http.Error(w, "Invalid forwarding URL: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.datastore.SetForwardingAddress(username, address); err != nil {
			log.Printf("failed to save forwarding address for user %s: %v", username, err)
			http.Error(w, "Failed to save forwarding address", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *defaultServer) forwardingAddressDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mustGetUsernameFromContext(r.Context())

		if err := s.datastore.DeleteForwardingAddress(username); err != nil {
			log.Printf("failed to delete forwarding address for user %s: %v", username, err)
			http.Error(w, "Failed to delete forwarding address", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
