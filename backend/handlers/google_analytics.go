package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/handlers/parse"
	"github.com/mtlynch/whatgotdone/backend/handlers/validate"
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

		users, err := s.datastore.Users()
		if err != nil {
			log.Printf("Failed to retrieve users from datastore: %v", err)
			http.Error(w, "Failed to retrieve pageviews", http.StatusInternalServerError)
			return
		}
		user, entryDate, ok := parseJournalEntryPath(path, users)
		if !ok {
			log.Printf("path is not a journal entry: %s", path)
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

		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
		}
	}
}

func (s defaultServer) wasEntryPublishedRecently(username types.Username, entryDate string) bool {
	entry, err := s.datastore.GetEntry(username, entryDate)
	if err != nil {
		return false
	}
	t, err := time.Parse(time.RFC3339, entry.LastModified)
	if err != nil {
		return false
	}
	return time.Now().Sub(t).Minutes() < 15
}

func (s *defaultServer) refreshGoogleAnalytics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.googleAnalyticsFetcher == nil {
			log.Print("Can't refresh Google Analytics because fetcher is not loaded")
			http.Error(w, "Google Analytics fetcher is not loaded", http.StatusInternalServerError)
			return
		}

		// Verify the request came from AppEngine so that external users can't
		// force the server to exceed Google Analytics rate limits.
		if !isAppEngineInternalRequest(r) {
			http.Error(w, "Refreshes of Google Analytics data must come from within AppEngine", http.StatusForbidden)
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
		var wg sync.WaitGroup
		for _, pvc := range pvcs {
			wg.Add(1)
			go func(pvc ga.PageViewCount) {
				defer wg.Done()
				if err := s.datastore.InsertPageViews(pvc.Path, pvc.Views); err != nil {
					log.Printf("failed to store pageviews in datastore %v: %v", pvc, err)
				}
			}(pvc)
		}
		wg.Wait()
		if err := json.NewEncoder(w).Encode(true); err != nil {
			panic(err)
		}
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
	users, err := s.datastore.Users()
	if err != nil {
		log.Printf("failed to retrieve user list for pageview filtering: %v", err)
		return filtered
	}
	for _, pvc := range pvcs {
		if _, _, ok := parseJournalEntryPath(pvc.Path, users); ok {
			filtered = append(filtered, pvc)
		}
	}
	return filtered
}

func parseJournalEntryPath(path string, users []types.Username) (types.Username, string, bool) {
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
	if !isUsernameInSlice(user, users) {
		return "", "", false
	}
	entryDate := pathParts[2]
	if !validate.EntryDate(entryDate) {
		return "", "", false
	}
	return user, entryDate, true
}

func isUsernameInSlice(u types.Username, uu []types.Username) bool {
	for _, x := range uu {
		if u == x {
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
