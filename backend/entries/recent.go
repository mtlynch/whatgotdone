package entries

import (
	"log"
	"sort"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type journalEntries []types.JournalEntry

func (r defaultReader) Recent(start, limit int) ([]types.JournalEntry, error) {
	// TODO: Filter by start date and min entry length.
	entries, err := r.store.ReadEntries(datastore.EntryFilter{})
	if err != nil {
		log.Printf("Failed to retrieve users: %s", err)
		return []RecentEntry{}, err
	}

	entriesRaw, err := r.recentEntriesFromAllUsers(users)
	if err != nil {
		return []RecentEntry{}, err
	}

	entries := recentEntries{}
	for _, entry := range entriesRaw {
		// Filter low-effort posts or test posts from the recent list.
		const minimumRelevantLength = 30
		if len(entry.Markdown) < minimumRelevantLength {
			continue
		}
		entries = append(entries, entry)
	}
	return sortAndSliceEntries(entries, start, limit), nil
}

func (r defaultReader) RecentFollowing(username types.Username, start, limit int) ([]RecentEntry, error) {
	following, err := r.store.Following(username)
	if err != nil {
		log.Printf("failed to retrieve user's follow list %s: %v", username, err)
		return []RecentEntry{}, err
	}

	var entries recentEntries
	for _, followedUsername := range following {
		entriesForUser, err := r.entriesFromUser(followedUsername)
		if err != nil {
			return recentEntries{}, err
		}
		entries = append(entries, entriesForUser...)
	}
	return sortAndSliceEntries(entries, start, limit), nil
}

func (r defaultReader) recentEntriesFromAllUsers(users []types.Username) ([]RecentEntry, error) {
	type result struct {
		entries []RecentEntry
		err     error
	}
	c := make(chan result)
	var wg sync.WaitGroup
	wg.Add(len(users))
	for _, username := range users {
		go func(u types.Username) {
			defer wg.Done()
			entriesForUser, err := r.entriesFromUser(u)
			c <- result{entriesForUser, err}
		}(username)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	entries := []RecentEntry{}
	var err error
	for res := range c {
		// Don't exit immediately because otherwise we'd leak the chan. Instead,
		//  save the first error we encounter.
		if err == nil && res.err != nil {
			err = res.err
		}
		entries = append(entries, res.entries...)
	}
	if err != nil {
		return []RecentEntry{}, err
	}

	return entries, nil
}

func (r defaultReader) entriesFromUser(username types.Username) (journalEntries, error) {
	entries := journalEntries{}
	journalEntries, err := r.store.ReadEntries(datastore.EntryFilter{
		ByUsers: []types.Username{username},
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
	return e[i].LastModified < e[j].LastModified
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
