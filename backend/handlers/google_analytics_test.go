package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func (ds mockDatastore) InsertPageViews(path string, pageViews int) error {
	return nil
}

func (ds mockDatastore) GetPageViews(path string) (int, error) {
	for _, pvc := range ds.pageViewCounts {
		if pvc.Path == path {
			return pvc.Views, nil
		}
	}
	return 0, errors.New("no pageview results found")
}

func (ds mockDatastore) GetEntry(username string, date string) (types.JournalEntry, error) {
	return types.JournalEntry{}, errors.New("not implemented")
}

func TestPageViewsGet(t *testing.T) {
	var pageViewsGetTests = []struct {
		path               string
		httpStatusExpected int
		viewsExpected      int
	}{
		{
			"/jimmy123/2020-01-17",
			http.StatusOK,
			5,
		},
		{
			"",
			http.StatusBadRequest,
			0,
		},
		{
			"/non.existent.user/2020-01-17",
			http.StatusForbidden,
			0,
		},
		{
			"/non.existent.user/2020-01-16",
			http.StatusForbidden,
			0,
		},
		{
			"/non.existent.user",
			http.StatusForbidden,
			0,
		},
		{
			"/",
			http.StatusForbidden,
			0,
		},
		{
			"/privacy-policy",
			http.StatusForbidden,
			0,
		},
	}

	ds := mockDatastore{
		users: []string{"jimmy123"},
		pageViewCounts: []ga.PageViewCount{
			ga.PageViewCount{"/jimmy123/2020-01-17", 5},
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	for _, tt := range pageViewsGetTests {
		req, err := http.NewRequest("GET", "/api/pageViews?path="+url.QueryEscape(tt.path), nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.httpStatusExpected {
			t.Fatalf("for input [%s], handler returned wrong status code: got %v want %v",
				tt.path, status, tt.httpStatusExpected)
		}

		if tt.httpStatusExpected == http.StatusOK {
			var response pageViewResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Response is not valid JSON: %v", w.Body.String())
			}
			if response.Views != tt.viewsExpected {
				t.Fatalf("for input [%s], unexpected view count: got %v want %v",
					tt.path, response.Views, tt.viewsExpected)
			}
		}
	}
}
