package redis

import (
	"errors"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetDraft returns an entry draft for the given user for the given date.
func (c client) GetDraft(username string, date string) (types.JournalEntry, error) {
	return types.JournalEntry{}, errors.New("not implemented")
}

// InsertDraft saves an entry draft to the datastore, overwriting any existing
// entry with the same name and username.
func (c client) InsertDraft(username string, j types.JournalEntry) error {
	return errors.New("not implemented")
}
