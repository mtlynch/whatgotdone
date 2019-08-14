package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mtlynch/whatgotdone/datastore"
	"github.com/mtlynch/whatgotdone/types"
)

func (s defaultServer) draftOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s defaultServer) draftGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to retrieve a draft entry", http.StatusForbidden)
			return
		}

		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		j, err := s.datastore.GetDraft(username, date)
		if _, ok := err.(datastore.DraftNotFoundError); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Failed to retrieve draft entry: %s", err)
			http.Error(w, "Failed to retrieve draft entry", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(j); err != nil {
			panic(err)
		}
	}
}

func (s defaultServer) draftPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to save a draft entry", http.StatusForbidden)
			return
		}

		type draftRequest struct {
			EntryContent string `json:"entryContent"`
		}

		type draftResponse struct {
			Ok bool `json:"ok"`
		}

		var t draftRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", err)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
		}

		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		j := types.JournalEntry{
			Date:         date,
			LastModified: time.Now().Format(time.RFC3339),
			Markdown:     t.EntryContent,
		}
		err = s.datastore.InsertDraft(username, j)
		if err != nil {
			log.Printf("Failed to update draft entry: %s", err)
			http.Error(w, "Failed to update draft entry", http.StatusInternalServerError)
			return
		}
		resp := draftResponse{
			Ok: true,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
