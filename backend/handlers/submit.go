package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s *defaultServer) submitOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}
	}
}

// submitPost handles HTTP POST requests for users to create new What Got Done
// updates. The updates can be new versions of previously published updates (in
// which case, we'll update the existing entries in the datastore) or a brand
// new update (in which case, we'll create new entries in the datastore).
func (s *defaultServer) submitPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to submit a journal entry", http.StatusForbidden)
			return
		}

		type submitRequest struct {
			Date         string `json:"date"`
			EntryContent string `json:"entryContent"`
		}

		var t submitRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", err)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
		}
		if !validateEntryDate(t.Date) {
			log.Printf("Invalid date: %s", t.Date)
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}

		j := types.JournalEntry{
			Date:         t.Date,
			LastModified: time.Now().Format(time.RFC3339),
			Markdown:     t.EntryContent,
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

		type submitResponse struct {
			Ok   bool   `json:"ok"`
			Path string `json:"path"`
		}
		resp := submitResponse{
			Ok:   true,
			Path: fmt.Sprintf("/%s/%s", username, t.Date),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}
