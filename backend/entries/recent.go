package entries

import (
	"log"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type journalEntries []types.JournalEntry

func (r defaultReader) Recent(start, limit int) ([]types.JournalEntry, error) {
	// TODO: Filter by start date.
	entries, err := r.store.ReadEntries(datastore.EntryFilter{
		// Filter low-effort posts.
		MinLength: 30,
		Offset:    int32(start),
		Limit:     int32(limit),
	})
	if err != nil {
		log.Printf("Failed to retrieve entries: %s", err)
		return journalEntries{}, err
	}

	return entries, nil
}

func (r defaultReader) RecentFollowing(username types.Username, start, limit int) ([]types.JournalEntry, error) {
	following, err := r.store.Following(username)
	if err != nil {
		log.Printf("failed to retrieve user's follow list %s: %v", username, err)
		return journalEntries{}, err
	}

	// TODO: Filter by start date.
	entries, err := r.store.ReadEntries(datastore.EntryFilter{
		ByUsers: following,
		Offset:  int32(start),
		Limit:   int32(limit),
	})
	if err != nil {
		log.Printf("Failed to retrieve entries: %s", err)
		return journalEntries{}, err
	}

	return entries, nil
}
