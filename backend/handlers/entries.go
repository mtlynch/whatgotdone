package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s *defaultServer) entriesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
			return
		}

		entries, err := s.datastore.GetEntries(username)
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", username), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

// entryPost handles HTTP POST requests for users to create new What Got
// Done updates. The updates can be new versions of previously published
// updates (in which case, we'll update the existing entries in the datastore)
// or a brand new update (in which case, we'll create new entries in the
// datastore).
func (s *defaultServer) entryPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to edit a journal entry", http.StatusForbidden)
			return
		}

		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		type entryRequest struct {
			EntryContent string `json:"entryContent"`
		}

		var t entryRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", err)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
		}

		j := types.JournalEntry{
			Date:             date,
			LastModifiedTime: time.Now(),
			Markdown:         t.EntryContent,
		}

		err = s.datastore.InsertDraft(username, j)
		if err != nil {
			log.Printf("Failed to update journal draft entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
			return
		}
		err = s.datastore.InsertEntry(username, j)
		if err != nil {
			log.Printf("Failed to insert journal entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
			return
		}

		type entryResponse struct {
			Ok   bool   `json:"ok"`
			Path string `json:"path"`
		}
		resp := entryResponse{
			Ok:   true,
			Path: fmt.Sprintf("/%s/%s", username, date),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
