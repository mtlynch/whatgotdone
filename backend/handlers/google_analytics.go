package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/handlers/validate"
)

type pageViewResponse struct {
	Path  string `json:"path"`
	Views int    `json:"views"`
}

func (s *defaultServer) pageViewsOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s defaultServer) pageViewsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			log.Print("Request is missing path query parameter")
			http.Error(w, "Request is missing path query parameter", http.StatusBadRequest)
			return
		}

		users, err := s.datastore.Users()
		if err != nil {
			log.Printf("Failed to retrieve users from datastore: %v", err)
			http.Error(w, "Failed to retrieve pageviews", http.StatusInternalServerError)
		}
		if !isPathForJournalEntry(path, users) {
			log.Printf("path is not a journal entry: %s", path)
			http.Error(w, "path parameter must specify a journal entry", http.StatusForbidden)
			return
		}

		views, err := s.datastore.GetPageViews(path)
		if _, ok := err.(datastore.PageViewsNotFoundError); ok {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("failed to retrieve pageviews from datastore for path %s: %v", path, err)
			http.Error(w, fmt.Sprintf("Failed to retrieve pageviews for path %s", path), http.StatusInternalServerError)
		}

		response := pageViewResponse{
			Path:  path,
			Views: views,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) refreshGoogleAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.googleAnalyticsFetcher == nil {
			http.Error(w, "Google Analytics fetcher is not loaded", http.StatusInternalServerError)
			return
		}
		// Verify the request came from AppEngine so that external users can't
		// force the server to exceed Google Analytics rate limits.
		if !isAppEngineInternalRequest(r) {
			http.Error(w, "Refreshes of Google Analytics data must come from within AppEngine", http.StatusForbidden)
			return
		}

		pvcs, err := (*s.googleAnalyticsFetcher).PageViewsByPath("2019-01-01", "today")
		if err != nil {
			log.Printf("failed to refresh Google Analytics data: %v", err)
			http.Error(w, "Failed to refresh Google Analytics data", http.StatusInternalServerError)
			return
		}
		pvcs = coalescePageViews(pvcs)
		pvcs = s.filterNonEntries(pvcs)
		for _, pvc := range pvcs {
			if err := s.datastore.InsertPageViews(pvc.Path, pvc.Views); err != nil {
				log.Printf("failed to store pageviews in datastore %v: %v", pvc, err)
			}
		}
		if err := json.NewEncoder(w).Encode(true); err != nil {
			panic(err)
		}
	}
}

func coalescePageViews(pvcs []ga.PageViewCount) []ga.PageViewCount {
	totals := map[string]int{}
	coalesced := []ga.PageViewCount{}
	for _, pvc := range pvcs {
		u, err := url.Parse(pvc.Path)
		if err != nil {
			panic(err)
		}
		if _, ok := totals[u.EscapedPath()]; ok {
			totals[u.EscapedPath()] += pvc.Views
		} else {
			totals[u.EscapedPath()] = pvc.Views
		}
	}
	for p, c := range totals {
		coalesced = append(coalesced, ga.PageViewCount{p, c})
	}
	return coalesced
}

func (s defaultServer) filterNonEntries(pvcs []ga.PageViewCount) []ga.PageViewCount {
	filtered := []ga.PageViewCount{}
	users, err := s.datastore.Users()
	if err != nil {
		return filtered
	}
	for _, pvc := range pvcs {
		if isPathForJournalEntry(pvc.Path, users) {
			filtered = append(filtered, pvc)
		}
	}
	return filtered
}

func isPathForJournalEntry(path string, users []string) bool {
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 3 {
		return false
	}
	// Path should start with a forward slash, so the first part is empty.
	if pathParts[0] != "" {
		return false
	}

	user := pathParts[1]
	if !isStringInSlice(user, users) {
		return false
	}
	entryDate := pathParts[2]
	if !validate.EntryDate(entryDate) {
		return false
	}
	return true
}

func isStringInSlice(s string, ss []string) bool {
	for _, x := range ss {
		if s == x {
			return true
		}
	}
	return false
}

func isAppEngineInternalRequest(r *http.Request) bool {
	cronHeader := r.Header.Get("X-Appengine-Cron")
	if cronHeader != "true" {
		log.Printf("X-Appengine-Cron=[%v]", cronHeader)
		return false
	}
	return true
}
