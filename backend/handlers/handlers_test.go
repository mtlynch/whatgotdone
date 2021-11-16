package handlers

import (
	"errors"
	"os"
	"path"
	"sync"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type mockDatastore struct {
	journalEntries []types.JournalEntry
	journalDrafts  []types.JournalEntry
	users          []types.Username
	reactions      []types.Reaction
	pageViewCounts []ga.PageViewCount
	userProfile    types.UserProfile
	mu             sync.Mutex
}

func (ds *mockDatastore) ReadEntries(filter datastore.EntryFilter) ([]types.JournalEntry, error) {
	return ds.journalEntries, nil
}

func (ds *mockDatastore) GetDraft(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	if len(ds.journalDrafts) > 0 {
		return ds.journalDrafts[0], nil
	}
	return types.JournalEntry{}, datastore.DraftNotFoundError{
		Username: username,
		Date:     date,
	}
}

func (ds *mockDatastore) InsertEntry(username types.Username, j types.JournalEntry) error {
	return nil
}

func (ds *mockDatastore) InsertDraft(username types.Username, j types.JournalEntry) error {
	return nil
}

func (ds *mockDatastore) InsertFollow(leader, follower types.Username) error {
	return errors.New("not implemented")
}

func (ds *mockDatastore) DeleteFollow(leader, follower types.Username) error {
	return errors.New("not implemented")
}

func (ds *mockDatastore) Following(follower types.Username) ([]types.Username, error) {
	return []types.Username{}, errors.New("not implemented")
}

func (ds *mockDatastore) GetPreferences(username types.Username) (types.Preferences, error) {
	return types.Preferences{}, datastore.PreferencesNotFoundError{
		Username: username,
	}
}

func (ds *mockDatastore) SetPreferences(username types.Username, prefs types.Preferences) error {
	return errors.New("not implemented")
}

func (ds *mockDatastore) Close() error {
	return nil
}

type mockAuthenticator struct {
	tokensToUsers map[string]types.Username
}

func (a mockAuthenticator) UserFromAuthToken(authToken string) (types.Username, error) {
	for k, v := range a.tokensToUsers {
		if k == authToken {
			return v, nil
		}
	}
	return "", errors.New("mock token not found")
}

func init() {
	// The handler uses relative paths to find the template file. Switch to the
	// app's root directory so that the relative paths work.
	if err := os.Chdir("../../"); err != nil {
		panic(err)
	}
	frontendIndexPath := path.Join(frontendRootDir, frontendIndexFilename)

	// Ensure that the frontend/dist/index.html exists. The handler functions
	// need it, even if it's empty.
	if _, err := os.Stat(frontendIndexPath); os.IsNotExist(err) {
		// Ensure that the frontend/dist folder exists.
		if err = os.MkdirAll(frontendRootDir, os.ModePerm); err != nil {
			panic(err)
		}
		// Create frontend/dist/index.html.
		if _, err := os.Create(frontendIndexPath); err != nil {
			panic(err)
		}
	}

}
