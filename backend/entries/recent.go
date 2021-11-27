package entries

import (
	"log"
	"sort"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type journalEntries []types.JournalEntry

func (r defaultReader) Recent(start, limit int) ([]types.JournalEntry, error) {
	// TODO: Filter by start date.
	entries, err := r.store.ReadEntries(datastore.EntryFilter{
		// Filter low-effort posts.
		MinLength: 30,
	})
	if err != nil {
		log.Printf("Failed to retrieve entries: %s", err)
		return journalEntries{}, err
	}

	return sortAndSliceEntries(entries, start, limit), nil
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
	})
	if err != nil {
		log.Printf("Failed to retrieve entries: %s", err)
		return journalEntries{}, err
	}

	return sortAndSliceEntries(entries, start, limit), nil
}

// TODO: Reimplement this in SQL.
func sortAndSliceEntries(entries journalEntries, start, limit int) journalEntries {
	sorted := make(journalEntries, len(entries))
	copy(sorted, entries)

	sort.Sort(sorted)
	// Reverse the order of entries.
	for i := len(sorted)/2 - 1; i >= 0; i-- {
		opp := len(sorted) - 1 - i
		sorted[i], sorted[opp] = sorted[opp], sorted[i]
	}

	start = min(len(sorted), start)
	end := min(len(sorted), start+limit)
	return sorted[start:end]
}

func (e journalEntries) Len() int {
	return len(e)
}

func (e journalEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e journalEntries) Less(i, j int) bool {
	if e[i].Date < e[j].Date {
		return true
	}
	if e[i].Date > e[j].Date {
		return false
	}
	return e[i].LastModified.Before(e[j].LastModified)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
