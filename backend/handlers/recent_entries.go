package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/mtlynch/whatgotdone/backend/entries"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type entryPublic struct {
	Author   types.Username `json:"author"`
	Date     string         `json:"date"`
	Markdown string         `json:"markdown"`
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

		entries, err := s.entriesReader.Recent(start, limit)
		if err != nil {
			log.Printf("Failed to retrieve recent entries: %v", err)
			http.Error(w, "Failed to retrieve recent entries", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(recentEntriesToPublicEntries(entries)); err != nil {
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

		entries, err := s.entriesReader.RecentFollowing(username, start, limit)
		if err != nil {
			log.Printf("Failed to retrieve recent entries from users %s is following: %v", username, err)
			http.Error(w, "Failed to retrieve recent entries from followed users", http.StatusInternalServerError)
			return
		}

		type response struct {
			Entries entriesPublic `json:"entries"`
		}
		resp := response{
			Entries: recentEntriesToPublicEntries(entries),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
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

func recentEntriesToPublicEntries(entries []entries.RecentEntry) entriesPublic {
	p := entriesPublic{}
	for _, entry := range entries {
		p = append(p, entryPublic{
			Author:   entry.Author,
			Date:     entry.Date,
			Markdown: entry.Markdown,
		})
	}
	return p
}
