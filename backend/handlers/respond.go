package handlers

import (
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, status int, data interface{}) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondOK(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusOK, data)
}
