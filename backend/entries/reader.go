package entries

import (
	"github.com/mtlynch/whatgotdone/backend/types"
)

type Reader interface {
	Recent(start, limit int) ([]RecentEntry, error)
	RecentFollowing(username string, start, limit int) ([]RecentEntry, error)
}

type entrystore interface {
	// Users returns all the users who have published entries.
	Users() ([]string, error)
	// GetEntries returns all published entries for the given user.
	GetEntries(username string) ([]types.JournalEntry, error)
	// GetReactions retrieves reader reactions associated with a published entry.
	// Followers returns all the users the specified user is following.
	Following(follower string) ([]string, error)
}

type defaultReader struct {
	store entrystore
}

func NewReader(store entrystore) Reader {
	return defaultReader{
		store: store,
	}
}
