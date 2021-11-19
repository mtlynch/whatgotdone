package handlers

import (
	"fmt"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type exportedEntry struct {
	Date         types.EntryDate `json:"date"`
	Markdown     string          `json:"markdown"`
	LastModified string          `json:"lastModified"`
}

func (s defaultServer) exportGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to retrieve a draft entry", http.StatusForbidden)
			return
		}

		entries, err := s.datastore.GetEntries(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve entries: %s", err), http.StatusInternalServerError)
			return
		}

		respondOK(w, struct {
			Entries []exportedEntry `json:"entries"`
		}{
			Entries: entriesToExportedEntries(entries, username),
		})
	}
}

func entriesToExportedEntries(entries []types.JournalEntry, author types.Username) []exportedEntry {
	p := []exportedEntry{}
	for _, entry := range entries {
		p = append(p, exportedEntry{
			Date:         entry.Date,
			Markdown:     entry.Markdown,
			LastModified: entry.LastModified,
		})
	}
	return p
}
