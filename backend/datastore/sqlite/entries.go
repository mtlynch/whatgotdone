package sqlite

import "github.com/mtlynch/whatgotdone/backend/types"

func (d db) GetEntry(username types.Username, date types.EntryDate) (types.JournalEntry, error) {
	return types.JournalEntry{}, notImplementedError
}

// InsertEntry saves an entry to the datastore, overwriting any existing entry
// with the same name and username.
func (d db) InsertEntry(username types.Username, j types.JournalEntry) error {
	return notImplementedError
}

// GetEntries returns all published entries for the given user.
func (d db) GetEntries(username types.Username) ([]types.JournalEntry, error) {
	return []types.JournalEntry{}, notImplementedError
}
