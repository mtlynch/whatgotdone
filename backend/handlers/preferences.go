package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (s defaultServer) preferencesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to retrieve your preferences", http.StatusForbidden)
			return
		}

		prefs, err := s.datastore.GetPreferences(username)
		if _, ok := err.(datastore.PreferencesNotFoundError); ok {
			http.Error(w, "No user preferences found", http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Failed to retrieve user preferences for %s: %v", username, err)
			http.Error(w, "Failed to retrieve preferences", http.StatusInternalServerError)
			return
		}

		respondOK(w, struct {
			EntryTemplate string `json:"entryTemplate"`
		}{
			EntryTemplate: prefs.EntryTemplate,
		})
	}
}

func (s defaultServer) preferencesPut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to update your preferences", http.StatusForbidden)
			return
		}
		prefs, err := preferencesFromRequest(r)
		if err != nil {
			http.Error(w, "Invalid preferences update request", http.StatusBadRequest)
			return
		}

		err = s.datastore.SetPreferences(username, prefs)
		if err != nil {
			log.Printf("failed to save updated preferences for user %s: %v", username, err)
			http.Error(w, "Failed to save preferences", http.StatusInternalServerError)
			return
		}
	}
}

func preferencesFromRequest(r *http.Request) (types.Preferences, error) {
	type preferencesRequest struct {
		EntryTemplate string `json:"entryTemplate"`
	}
	var pr preferencesRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pr)
	if err != nil {
		return types.Preferences{}, err
	}
	return types.Preferences{
		EntryTemplate: strings.TrimSpace(pr.EntryTemplate),
	}, nil
}
