package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type entryPublic struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	// Skip JSON serialization for lastModified as clients don't need this field,
	// but we need it internally for sorting lists of entries.
	lastModified string
	Markdown     string `json:"markdown"`
}

type entriesPublic []entryPublic

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

		entriesFull, err := s.entriesReader.Recent()
		if err != nil {
			log.Printf("Failed to retrieve recent entries: %v", err)
			http.Error(w, "Failed to retrieve recent entries", http.StatusInternalServerError)
			return
		}

		entries := entriesPublic{}
		for _, entry := range userEntries {
			// Filter low-effort posts or test posts from the recent list.
			const minimumRelevantLength = 30
			if len(entry.Markdown) < minimumRelevantLength {
				continue
			}
			entries = append(entries, entryPublic{
				Author:       username,
				Date:         entry.Date,
				lastModified: entry.LastModified,
				Markdown:     entry.Markdown,
			})
		}

		if err := json.NewEncoder(w).Encode(sortAndSliceEntries(entries, start, limit)); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) entriesFollowingGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to retrieve your personalized feed", http.StatusForbidden)
			return
		}
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

		following, err := s.datastore.Following(username)
		if err != nil {
			log.Printf("failed to retrieve user's follow list %s: %v", username, err)
			http.Error(w, "Failed to retrieve your personalized feed", http.StatusInternalServerError)
			return
		}

		var entries entriesPublic
		for _, followedUsername := range following {
			userEntries, err := s.datastore.GetEntries(followedUsername)
			if err != nil {
				log.Printf("Failed to retrieve entries for user %s: %v", followedUsername, err)
				http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", followedUsername), http.StatusInternalServerError)
				return
			}
			for _, entry := range userEntries {
				entries = append(entries, entryPublic{
					Author:       followedUsername,
					Date:         entry.Date,
					lastModified: entry.LastModified,
					Markdown:     entry.Markdown,
				})
			}
		}

		type response struct {
			Entries []entryPublic `json:"entries"`
		}
		resp := response{
			Entries: sortAndSliceEntries(entries, start, limit),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func sortAndSliceEntries(entries entriesPublic, start, limit int) entriesPublic {
	sorted := make(entriesPublic, len(entries))
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

func (e entriesPublic) Len() int {
	return len(e)
}

func (e entriesPublic) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e entriesPublic) Less(i, j int) bool {
	if e[i].Date < e[j].Date {
		return true
	}
	if e[i].Date > e[j].Date {
		return false
	}
	return e[i].lastModified < e[j].lastModified
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
