package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/handlers/parse"
	"github.com/mtlynch/whatgotdone/backend/types"
)

type pageViewResponse struct {
	Path  string `json:"path"`
	Views int    `json:"views"`
}

func (s defaultServer) pageViewsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Query().Get("path")
		if path == "" {
			log.Print("Request is missing path query parameter")
			http.Error(w, "Request is missing path query parameter", http.StatusBadRequest)
			return
		}

		user, entryDate, ok := parseJournalEntryPath(path)
		if !ok {
			log.Printf("path is not a journal entry: %s", path)
			http.Error(w, "path parameter must specify a journal entry", http.StatusForbidden)
			return
		}

		// Make sure the requested path actually has a journal entry associated with
		// it.
		_, err := s.datastore.GetEntry(user, entryDate)
		if err != nil {
			log.Printf("failed to find entry associated with requested pageview: %s, %s -> %v", user, entryDate, err)
			http.Error(w, "path parameter must specify a journal entry", http.StatusForbidden)
			return
		}

		views, err := s.datastore.GetPageViews(path)
		if _, ok := err.(datastore.PageViewsNotFoundError); ok {
			if s.wasEntryPublishedRecently(user, entryDate) {
				views = 1
				err = nil
			} else {
				log.Printf("No pageviews found for %s", path)
				http.Error(w, "Path has no pageview data", http.StatusNotFound)
				return
			}
		} else if err != nil {
			log.Printf("failed to retrieve pageviews from datastore for path %s: %v", path, err)
			http.Error(w, fmt.Sprintf("Failed to retrieve pageviews for path %s", path), http.StatusInternalServerError)
			return
		}

		response := pageViewResponse{
			Path:  path,
			Views: views,
		}

		respondOK(w, response)
	}
}

func (s defaultServer) wasEntryPublishedRecently(username types.Username, entryDate types.EntryDate) bool {
	entry, err := s.datastore.GetEntry(username, entryDate)
	if err != nil {
		return false
	}
	return time.Now().Sub(entry.LastModified).Minutes() < 15
}

func (s *defaultServer) refreshGoogleAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.googleAnalyticsFetcher == nil {
			log.Print("Can't refresh Google Analytics because fetcher is not loaded")
			http.Error(w, "Google Analytics fetcher is not loaded", http.StatusInternalServerError)
			return
		}

		pvcs, err := s.googleAnalyticsFetcher.PageViewsByPath("2019-01-01", "today")
		if err != nil {
			log.Printf("failed to refresh Google Analytics data: %v", err)
			http.Error(w, "Failed to refresh Google Analytics data", http.StatusInternalServerError)
			return
		}
		pvcs = coalescePageViews(pvcs)
		pvcs = s.filterNonEntries(pvcs)

		err = s.datastore.InsertPageViews(pvcs)
		if err != nil {
			log.Printf("failed to store Google Analytics data: %v", err)
			http.Error(w, "Failed to save Google Analytics data", http.StatusInternalServerError)
			return
		}

		respondOK(w, true)
	}
}

// coalescePageViews aggregates together path strings with the same path part
// but different query strings.
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
	for p, v := range totals {
		coalesced = append(coalesced, ga.PageViewCount{
			Path:  p,
			Views: v,
		})
	}
	return coalesced
}

// filterNonEntries removes page counts for paths that are not user entries.
func (s defaultServer) filterNonEntries(pvcs []ga.PageViewCount) []ga.PageViewCount {
	filtered := []ga.PageViewCount{}
	for _, pvc := range pvcs {
		if _, _, ok := parseJournalEntryPath(pvc.Path); ok {
			filtered = append(filtered, pvc)
		}
	}
	return filtered
}

func parseJournalEntryPath(path string) (types.Username, types.EntryDate, bool) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 3 {
		return "", "", false
	}
	// Path should start with a forward slash, so the first part is empty.
	if pathParts[0] != "" {
		return "", "", false
	}

	user, err := parse.Username(pathParts[1])
	if err != nil {
		return "", "", false
	}
	entryDate, err := parse.EntryDate(pathParts[2])
	if err != nil {
		return "", "", false
	}
	return user, entryDate, true
}
