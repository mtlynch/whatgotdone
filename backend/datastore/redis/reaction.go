package redis

import (
	"errors"

	"github.com/mtlynch/whatgotdone/backend/types"
)

// GetReactions retrieves reader reactions associated with a published entry.
func (c client) GetReactions(entryAuthor string, entryDate string) ([]types.Reaction, error) {
	return []types.Reaction{}, errors.New("not implemented")
}

// AddReaction saves a reader reaction associated with a published entry,
// overwriting any existing reaction.
func (c client) AddReaction(entryAuthor string, entryDate string, reaction types.Reaction) error {
	return errors.New("not implemented")
}
