package sqlite

import "github.com/mtlynch/whatgotdone/backend/types"

// GetReactions retrieves reader reactions associated with a published entry.
func (d db) GetReactions(entryAuthor types.Username, entryDate types.EntryDate) ([]types.Reaction, error) {
	return []types.Reaction{}, notImplementedError
}

// AddReaction saves a reader reaction associated with a published entry,
// overwriting any existing reaction.
func (d db) AddReaction(entryAuthor types.Username, entryDate types.EntryDate, reaction types.Reaction) error {
	return notImplementedError
}
