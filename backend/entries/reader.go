package entries

import (
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Reader reads journal entries.
type Reader interface {
	// Recent returns the recent entries in the store.
	Recent(start, limit int) ([]RecentEntry, error)
	// RecentFollowing returns recent entries from among users that the specified
	// user is following.
	RecentFollowing(username string, start, limit int) ([]RecentEntry, error)
}

// EntryStore stores information related to journal entries.
type EntryStore interface {
	// Users returns all the users who have published entries.
	Users() ([]string, error)
	// GetEntries returns all published entries for the given user.
	GetEntries(username string) ([]types.JournalEntry, error)
	// GetReactions retrieves reader reactions associated with a published entry.
	// Followers returns all the users the specified user is following.
	Following(follower string) ([]string, error)
}

type defaultReader struct {
	store EntryStore
}

// NewReader creates a new entries.Reader.
func NewReader(store EntryStore) Reader {
	return defaultReader{
		store: store,
	}
}
