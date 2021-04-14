package entries

import (
	"log"

	"github.com/mtlynch/whatgotdone/backend/types"
)

type RecentEntry struct {
	Author       string
	Date         string
	LastModified string
	Markdown     string
}

func (r defaultReader) Recent() ([]types.JournalEntry, error) {
	users, err := r.datastore.Users()
	if err != nil {
		log.Printf("Failed to retrieve users: %s", err)
		return []types.JournalEntry{}, err
	}

	entries := []types.JournalEntry{}
	for _, username := range users {
		userEntries, err := r.datastore.GetEntries(username)
		if err != nil {
			log.Printf("Failed to retrieve entries for user %s: %s", username, err)
			return []types.JournalEntry{}, err
		}
		for _, entry := range userEntries {
			// Filter low-effort posts or test posts from the recent list.
			const minimumRelevantLength = 30
			if len(entry.Markdown) < minimumRelevantLength {
				continue
			}
			entries = append(entries, entry)
		}
	}
	return entries, nil
}
