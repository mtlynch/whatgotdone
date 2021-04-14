package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

		entriesFull, err := s.entriesReader.Recent(start, limit)
		if err != nil {
			log.Printf("Failed to retrieve recent entries: %v", err)
			http.Error(w, "Failed to retrieve recent entries", http.StatusInternalServerError)
			return
		}

		entries := entriesPublic{}
		for _, entry := range entriesFull {
			entries = append(entries, entryPublic{
				Author:   entry.Author,
				Date:     entry.Date,
				Markdown: entry.Markdown,
			})
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
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
			//log.Printf("Failed to retrieve recent entries: %v", err)
			//http.Error(w, "Failed to retrieve recent entries", http.StatusInternalServerError)
			return
		}

		type response struct {
			Entries []entryPublic `json:"entries"`
		}
		resp := response{
			Entries: entries,
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
