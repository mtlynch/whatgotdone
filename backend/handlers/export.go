package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
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

	log.Printf("exporting(%s): unpublished drafts", username)
	drafts, err := s.exportUserDrafts(username)
	if err != nil {
		return exportedUserData{}, err
	}

	log.Printf("exporting(%s): published entries", username)
	entries, err := s.datastore.GetEntries(username)
	if err != nil {
		return exportedUserData{}, err
	}

	log.Printf("exporting(%s): preferences", username)
	prefs, err := s.datastore.GetPreferences(username)
	if _, ok := err.(datastore.PreferencesNotFoundError); ok {
		prefs = types.Preferences{}
	} else if err != nil {
		return exportedUserData{}, err
	}

	log.Printf("exporting(%s): user profile", username)
	profile, err := s.datastore.GetUserProfile(username)
	if _, ok := err.(datastore.UserProfileNotFoundError); ok {
		profile = types.UserProfile{}
	} else if err != nil {
		return exportedUserData{}, err
	}

	log.Printf("exporting(%s): followed users", username)
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
	// Retrieve all the user's drafts by checking every possible draft date for an
	// entry. This is inefficient, and we could optimize/parallelize this, but
	// exporting isn't a very common or performance-sensitive code path.

	type result struct {
		draft types.JournalEntry
		err   error
	}
	c := make(chan result)
	var wg sync.WaitGroup

	// 2019-03-29 is the first ever post on What Got Done.
	currentDate := time.Date(2019, time.March, 29, 0, 0, 0, 0, time.UTC)
	for {
		if currentDate.After(dates.ThisFriday()) {
			break
		}
		wg.Add(1)
		go func(date types.EntryDate) {
			defer wg.Done()
			draft, err := s.datastore.GetDraft(username, date)
			c <- result{draft, err}
		}(types.EntryDate(currentDate.Format("2006-01-02")))

		// Increment to next Friday.
		currentDate = currentDate.AddDate(0, 0, 7)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	drafts := []types.JournalEntry{}
	var err error
	for res := range c {
		if res.err == nil {
			drafts = append(drafts, res.draft)
			continue
		}
		if _, ok := err.(datastore.DraftNotFoundError); ok {
			continue
		}
		// Don't exit immediately because otherwise we'd leak the chan. Instead,
		//  save the first error we encounter.
		if err == nil && res.err != nil {
			err = res.err
		}
	}

	// Sort drafts in ascending order of date.
	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].Date < drafts[j].Date
	})
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
