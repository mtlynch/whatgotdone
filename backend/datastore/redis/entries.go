package redis

import (
	"errors"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetEntries returns all published entries for the given user.
func (c client) GetEntries(username string) ([]types.JournalEntry, error) {
	return []types.JournalEntry{}, errors.New("not implemented")
}

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (c client) InsertEntry(username string, j types.JournalEntry) error {
	return errors.New("not implemented")
}
