package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
)

func (s *defaultServer) pageViewsOptions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (s defaultServer) pageViewsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// if only one expected
		path := r.URL.Query().Get("path")
		if path == "" {
			log.Printf("Request is missing path query parameter")
			http.Error(w, "Request is missing path query parameter", http.StatusBadRequest)
			return
		}

		// TODO: Validate the path.

		// TODO: Handle missing entry in datastore.
		views, err := s.datastore.GetPageViews(path)
		if err != nil {
			log.Printf("failed to retrieve pageviews from datastore for path %s: %v", path, err)
			http.Error(w, fmt.Sprintf("Failed to retrieve pageviews for path %s", path), http.StatusInternalServerError)
		}

		type pageViewResponse struct {
			Path  string `json:"path"`
			Views int    `json:"views"`
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
		// TODO: Verify the request came from AppEngine.
		pvcs, err := (*s.googleAnalyticsFetcher).PageViewsByPath("2019-01-01", "today")
		if err != nil {
			log.Printf("failed to refresh Google Analytics data: %v", err)
			http.Error(w, "Failed to refresh Google Analytics data", http.StatusInternalServerError)
			return
		}
		for _, pvc := range coalescePageViews(pvcs) {
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
