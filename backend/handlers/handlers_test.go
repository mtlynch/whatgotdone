package handlers

import (
	"errors"
	"os"
	"path"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type mockDatastore struct {
	journalEntries []types.JournalEntry
	journalDrafts  []types.JournalEntry
	users          []string
	reactions      []types.Reaction
}

func (ds mockDatastore) Users() ([]string, error) {
	return ds.users, nil
}

func (ds mockDatastore) AllEntries(username string) ([]types.JournalEntry, error) {
	return ds.journalEntries, nil
}

func (ds mockDatastore) GetDraft(username string, date string) (types.JournalEntry, error) {
	if len(ds.journalDrafts) > 0 {
		return ds.journalDrafts[0], nil
	}
	return types.JournalEntry{}, datastore.DraftNotFoundError{
		Username: username,
		Date:     date,
	}
}

func (ds mockDatastore) InsertEntry(username string, j types.JournalEntry) error {
	return nil
}

func (ds mockDatastore) InsertDraft(username string, j types.JournalEntry) error {
	return nil
}

func (ds mockDatastore) Close() error {
	return nil
}

type mockAuthenticator struct {
	tokensToUsers map[string]string
}

func (a mockAuthenticator) UserFromAuthToken(authToken string) (string, error) {
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
