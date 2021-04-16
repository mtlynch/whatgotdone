package entries

import (
	"log"
	"sort"
)

// RecentEntry stores data about a journal entry.
type RecentEntry struct {
	Author       string
	Date         string
	LastModified string
	Markdown     string
}

type recentEntries []RecentEntry

func (r defaultReader) Recent(start, limit int) ([]RecentEntry, error) {
	users, err := r.store.Users()
	if err != nil {
		log.Printf("Failed to retrieve users: %s", err)
		return []RecentEntry{}, err
	}

	entries := []RecentEntry{}
	for _, username := range users {
		userEntries, err := r.store.GetEntries(username)
		if err != nil {
			log.Printf("Failed to retrieve entries for user %s: %s", username, err)
			return []RecentEntry{}, err
		}
		for _, entry := range userEntries {
			// Filter low-effort posts or test posts from the recent list.
			const minimumRelevantLength = 30
			if len(entry.Markdown) < minimumRelevantLength {
				continue
			}
			entries = append(entries, RecentEntry{
				Author:       username,
				Date:         entry.Date,
				LastModified: entry.LastModified,
				Markdown:     entry.Markdown,
			})
		}
	}
	return sortAndSliceEntries(entries, start, limit), nil
}

func (r defaultReader) RecentFollowing(username string, start, limit int) ([]RecentEntry, error) {
	following, err := r.store.Following(username)
	if err != nil {
		log.Printf("failed to retrieve user's follow list %s: %v", username, err)
		return []RecentEntry{}, err
	}

	var entries recentEntries
	for _, followedUsername := range following {
		userEntries, err := r.store.GetEntries(followedUsername)
		if err != nil {
			log.Printf("Failed to retrieve entries for user %s: %v", followedUsername, err)
			return []RecentEntry{}, err
		}
		for _, entry := range userEntries {
			entries = append(entries, RecentEntry{
				Author:       followedUsername,
				Date:         entry.Date,
				LastModified: entry.LastModified,
				Markdown:     entry.Markdown,
			})
		}
	}
	return sortAndSliceEntries(entries, start, limit), nil
}

func sortAndSliceEntries(entries recentEntries, start, limit int) recentEntries {
	sorted := make(recentEntries, len(entries))
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

func (e recentEntries) Len() int {
	return len(e)
}

func (e recentEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e recentEntries) Less(i, j int) bool {
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
