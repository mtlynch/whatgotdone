package handlers

import (
	"encoding/json"
	"net/http"
)

func encodeBody(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func respond(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	if data != nil {
		encodeBody(w, data)
	}
}

func respondOK(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusOK, data)
}
