package handlers

import (
	"fmt"
	"net/http"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type exportedEntry struct {
	Date         types.EntryDate `json:"date"`
	Markdown     string          `json:"markdown"`
	LastModified string          `json:"lastModified"`
}

type exportedPreferences struct {
	EntryTemplate string `json:"entryTemplate"`
}

func (s defaultServer) exportGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to retrieve a draft entry", http.StatusForbidden)
			return
		}

		// TODO: Get drafts

		entries, err := s.datastore.GetEntries(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve published entries: %s", err), http.StatusInternalServerError)
			return
		}

		prefs, err := s.datastore.GetPreferences(username)
		if _, ok := err.(datastore.PreferencesNotFoundError); ok {
			prefs = types.Preferences{}
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve user preferences: %s", err), http.StatusInternalServerError)
			return
		}

		profile, err := s.datastore.GetUserProfile(username)
		if _, ok := err.(datastore.UserProfileNotFoundError); ok {
			profile = types.UserProfile{}
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve user profile: %s", err), http.StatusInternalServerError)
			return
		}

		following, err := s.datastore.Following(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve following users: %s", err), http.StatusInternalServerError)
			return
		}

		respondOK(w, struct {
			Entries     []exportedEntry     `json:"entries"`
			Drafts      []exportedEntry     `json:"drafts"`
			Following   []types.Username    `json:"following"`
			Profile     profilePublic       `json:"profile"`
			Preferences exportedPreferences `json:"preferences"`
		}{
			Entries:   entriesToExportedEntries(entries, username),
			Following: following,
			Profile:   profileToPublic(profile),
			Preferences: exportedPreferences{
				EntryTemplate: prefs.EntryTemplate,
			}})
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
