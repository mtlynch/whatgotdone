package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type entryPublic struct {
	Author   types.Username  `json:"author"`
	Date     types.EntryDate `json:"date"`
	Markdown string          `json:"markdown"`
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

		// TODO: Filter by start date.
		entries, err := s.datastore.ReadEntries(datastore.EntryFilter{
			// Filter low-effort posts.
			MinLength: 30,
			Offset:    int32(start),
			Limit:     int32(limit),
		})
		if err != nil {
			log.Printf("Failed to retrieve entries: %v", err)
			http.Error(w, fmt.Sprintf("Failed to read journal entries: %v", err), http.StatusInternalServerError)
			return
		}

		respondOK(w, entriesToPublicEntries(entries))
	}
}

func (s *defaultServer) entriesFollowingGet() http.HandlerFunc {
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

		username := mustGetUsernameFromContext(r.Context())
		following, err := s.datastore.Following(username)
		if err != nil {
			log.Printf("failed to retrieve user's follow list %s: %v", username, err)
			http.Error(w, "Failed to retrieve user's follow list", http.StatusInternalServerError)
			return
		}

		// TODO: Filter by start date.
		entries, err := s.datastore.ReadEntries(datastore.EntryFilter{
			ByUsers: following,
			Offset:  int32(start),
			Limit:   int32(limit),
		})
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, "Failed to retrieve followed entries", http.StatusInternalServerError)
		}

		respondOK(w, struct {
			Entries entriesPublic `json:"entries"`
		}{
			Entries: entriesToPublicEntries(entries),
		})
	}
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

func entriesToPublicEntries(entries []types.JournalEntry) entriesPublic {
	p := entriesPublic{}
	for _, entry := range entries {
		p = append(p, entryPublic{
			Author:   entry.Author,
			Date:     entry.Date,
			Markdown: string(entry.Markdown),
		})
	}
	return p
}
