package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
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

		entries, err := s.datastore.ReadEntries(datastore.EntryFilter{
			ByUsers: []types.Username{username},
		})
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", username), http.StatusInternalServerError)
			return
		}
		respondOK(w, entries)
	}
}

// entryPut handles HTTP PUT requests for users to create new What Got Done
// updates. The updates can be new versions of previously published updates (in
// which case, we'll update the existing entries in the datastore) or a brand
// new update (in which case, we'll create new entries in the datastore).
func (s *defaultServer) entryPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		c, err := entryContentFromRequest(r)
		if err != nil {
			log.Printf("Invalid entry request: %v", err)
			http.Error(w, fmt.Sprintf("Invalid entry request: %v", err), http.StatusBadRequest)
			return
		}

		j := types.JournalEntry{
			Date:     date,
			Markdown: c,
		}

		username := usernameFromContext(r.Context())

		// First update the latest draft entry.
		err = s.datastore.InsertDraft(username, j)
		if err != nil {
			log.Printf("Failed to update journal draft entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
			return
		}
		// Then, update the published version.
		err = s.datastore.InsertEntry(username, j)
		if err != nil {
			log.Printf("Failed to insert journal entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
			return
		}

		respondOK(w, struct {
			Path string `json:"path"`
		}{
			Path: fmt.Sprintf("/%s/%s", username, date),
		})
	}
}

func (s *defaultServer) entryDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		username := usernameFromContext(r.Context())
		err = s.datastore.DeleteEntry(username, date)
		if err != nil {
			log.Printf("Failed to delete journal entry: %s", err)
			http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
			return
		}
	}
}
