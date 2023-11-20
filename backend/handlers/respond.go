package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func respond(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatalf("failed to encode to JSON: %v", err)
		}
	}
}

func respondOK(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusOK, data)
}
