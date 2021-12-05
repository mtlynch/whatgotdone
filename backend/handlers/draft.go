package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) draftGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := usernameFromContext(r.Context())

		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		draftMarkdown, err := s.savedDraftOrEntryTemplate(username, date)
		if err != nil {
			log.Printf("Failed to retrieve draft entry: %s", err)
			http.Error(w, "Failed to retrieve draft entry", http.StatusInternalServerError)
			return
		}
		if draftMarkdown == "" {
			http.Error(w, "No draft found for this entry", http.StatusNotFound)
			return
		}

		respondOK(w, struct {
			Markdown string `json:"markdown"`
		}{
			Markdown: draftMarkdown,
		})
	}
}

func (s defaultServer) draftPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date, err := dateFromRequestPath(r)
		if err != nil {
			log.Printf("Invalid date: %s", date)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		type draftRequest struct {
			EntryContent string `json:"entryContent"`
		}

		var t draftRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", err)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
			return
		}

		username := usernameFromContext(r.Context())
		err = s.datastore.InsertDraft(username, types.JournalEntry{
			Date:     date,
			Markdown: t.EntryContent,
		})
		if err != nil {
			log.Printf("Failed to update draft entry: %s", err)
			http.Error(w, "Failed to update draft entry", http.StatusInternalServerError)
			return
		}
	}
}

func (s defaultServer) savedDraftOrEntryTemplate(username types.Username, date types.EntryDate) (string, error) {
	// First, check if there's a saved draft.
	d, err := s.datastore.GetDraft(username, date)
	if err == nil && d.Markdown != "" {
		return d.Markdown, nil
	}
	if _, ok := err.(datastore.DraftNotFoundError); ok {
		err = nil
	} else if err != nil {
		return "", err
	}

	// If there's no saved draft, try using the user's entry template.
	prefs, err := s.datastore.GetPreferences(username)
	if _, ok := err.(datastore.PreferencesNotFoundError); ok {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return prefs.EntryTemplate, nil
}
