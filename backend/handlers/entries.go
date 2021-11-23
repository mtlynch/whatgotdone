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
		respondOK(w, entries)
	}
}

// entryPut handles HTTP POST requests for users to create new What Got
// Done updates. The updates can be new versions of previously published
// updates (in which case, we'll update the existing entries in the datastore)
// or a brand new update (in which case, we'll create new entries in the
// datastore).
func (s *defaultServer) entryPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			Date:         date,
			LastModified: time.Now().Format(time.RFC3339),
			Markdown:     t.EntryContent,
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
