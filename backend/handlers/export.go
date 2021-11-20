package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/dates"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type (
	exportedEntry struct {
		Date         types.EntryDate `json:"date"`
		Markdown     string          `json:"markdown"`
		LastModified string          `json:"lastModified"`
	}

	exportedPreferences struct {
		EntryTemplate string `json:"entryTemplate"`
	}

	exportedUserData struct {
		Entries     []exportedEntry     `json:"entries"`
		Drafts      []exportedEntry     `json:"drafts"`
		Following   []types.Username    `json:"following"`
		Profile     profilePublic       `json:"profile"`
		Preferences exportedPreferences `json:"preferences"`
	}
)

func (s defaultServer) exportGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to export your data", http.StatusForbidden)
			return
		}

		d, err := s.exportUserData(username)
		if err != nil {
			log.Printf("Failed to export user data: %v", err)
			http.Error(w, fmt.Sprintf("Failed to export user data: %s", err), http.StatusInternalServerError)
			return
		}

		respondOK(w, d)
	}
}

func (s defaultServer) exportUserData(username types.Username) (exportedUserData, error) {
	log.Printf("starting export for %s", username)
	drafts, err := s.exportUserDrafts(username)
	if err != nil {
		return exportedUserData{}, err
	}

	entries, err := s.datastore.GetEntries(username)
	if err != nil {
		return exportedUserData{}, err
	}

	prefs, err := s.datastore.GetPreferences(username)
	if _, ok := err.(datastore.PreferencesNotFoundError); ok {
		prefs = types.Preferences{}
	} else if err != nil {
		return exportedUserData{}, err
	}

	profile, err := s.datastore.GetUserProfile(username)
	if _, ok := err.(datastore.UserProfileNotFoundError); ok {
		profile = types.UserProfile{}
	} else if err != nil {
		return exportedUserData{}, err
	}

	following, err := s.datastore.Following(username)
	if err != nil {
		return exportedUserData{}, err
	}

	log.Printf("finished export for %s", username)

	return exportedUserData{
		Entries:   entriesToExportedEntries(entries, username),
		Drafts:    entriesToExportedEntries(drafts, username),
		Following: following,
		Profile:   profileToPublic(profile),
		Preferences: exportedPreferences{
			EntryTemplate: prefs.EntryTemplate,
		},
	}, nil
}

func (s defaultServer) exportUserDrafts(username types.Username) ([]types.JournalEntry, error) {
	drafts := []types.JournalEntry{}

	// Retrieve all the user's drafts by checking every possible draft date for an
	// entry. This is inefficient, and we could optimize/parallelize this, but
	// exporting isn't a very common or performance-sensitive code path.

	// 2019-03-29 is the first ever post on What Got Done.
	currentDate := time.Date(2019, time.March, 29, 0, 0, 0, 0, time.UTC)
	for {
		if currentDate.After(dates.ThisFriday()) {
			break
		}
		draft, err := s.datastore.GetDraft(username, types.EntryDate(currentDate.Format("2006-01-02")))
		if _, ok := err.(datastore.DraftNotFoundError); ok {
			// Ignore not found errors
		} else if err != nil {
			return []types.JournalEntry{}, err
		}
		if err == nil {
			drafts = append(drafts, draft)
		}

		currentDate = currentDate.AddDate(0, 0, 7)
	}

	return drafts, nil
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
