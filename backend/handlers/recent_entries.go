package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type recentEntry struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	// Skip JSON serialization for lastModified as clients don't need this field.
	lastModified string
	Markdown     string `json:"markdown"`
}

func (s *defaultServer) recentEntriesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start, err := parseStart(r.URL.Query().Get("start"))
		if err != nil {
			http.Error(w, "Invalid start parameter", http.StatusBadRequest)
			return
		}
		limit, err := parseLimit(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}

		users, err := s.datastore.Users()
		if err != nil {
			log.Printf("Failed to retrieve users: %s", err)
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}

		entries := []recentEntry{}
		for _, username := range users {
			userEntries, err := s.datastore.GetEntries(username)
			if err != nil {
				log.Printf("Failed to retrieve entries for user %s: %s", username, err)
				http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
				return
			}
			for _, entry := range userEntries {
				// Filter low-effort posts or test posts from the recent list.
				const minimumRelevantLength = 30
				if len(entry.Markdown) < minimumRelevantLength {
					continue
				}
				entries = append(entries, recentEntry{
					Author:       username,
					Date:         entry.Date,
					lastModified: entry.LastModified,
					Markdown:     entry.Markdown,
				})
			}
		}

		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Date < entries[j].Date {
				return true
			}
			if entries[i].Date > entries[j].Date {
				return false
			}
			return entries[i].lastModified < entries[j].lastModified
		})

		// Reverse the order of entries.
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}

		start = min(len(entries), start)
		end := min(len(entries), start+limit)
		entries = entries[start:end]

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseStart(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, errors.New("start value can't be negative")
	}
	return i, nil
}

func parseLimit(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < 1 {
		return 0, errors.New("limit value must be positive")
	}
	return i, nil
}
