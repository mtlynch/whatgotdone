package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	ga "github.com/mtlynch/whatgotdone/backend/google_analytics"
	"github.com/mtlynch/whatgotdone/backend/types"
)

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

	ds := mock.MockDatastore{
		Usernames: []types.Username{"jimmy123"},
		JournalEntries: []types.JournalEntry{
			{
				Author: types.Username("jimmy123"),
				Date:   types.EntryDate("2020-01-17"),
			},
		},
		LastPageViewUpdate: time.Now().Add(time.Minute * -1),
	}
	ds.InsertPageViews([]ga.PageViewCount{
		{
			Path:  "/jimmy123/2020-01-17",
			Views: 5,
		},
	})
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
		s.Router().ServeHTTP(w, req)

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

type mockGoogleAnalyticsFetcher struct {
	PageViewCounts []ga.PageViewCount
}

func (f mockGoogleAnalyticsFetcher) PageViewsByPath(_, _ string) ([]ga.PageViewCount, error) {
	return f.PageViewCounts, nil
}

func TestPageViewsGetUpdatesPageViewsWhenRecordsAreStale(t *testing.T) {
	mf := mockGoogleAnalyticsFetcher{
		PageViewCounts: []ga.PageViewCount{
			{
				Path:  "/joe/2020-04-24",
				Views: 5,
			},
			{
				Path:  "/joe/2020-04-24?fbclid=dummy_facebook_referral_id",
				Views: 6,
			},
			{
				Path:  "/joe/2020-04-17",
				Views: 8,
			},
			{
				Path:  "/mary/2020-04-17",
				Views: 25,
			},
			{
				Path:  "/undefined/2020-04-17",
				Views: 100,
			},
			{
				Path:  "/entry/edit/2020-04-17",
				Views: 15,
			},
			{
				Path:  "/privacy-policy",
				Views: 2,
			},
		},
	}

	ds := mock.MockDatastore{
		Usernames: []types.Username{"joe", "mary"},
		JournalEntries: []types.JournalEntry{
			{
				Author: types.Username("joe"),
				Date:   types.EntryDate("2020-04-24"),
			},
		},
		LastPageViewUpdate:     time.Now().Add(time.Minute * -10),
		CallsToInsertPageViews: make(chan bool),
	}

	go func() {
		ds.InsertPageViews([]ga.PageViewCount{
			{
				Path:  "/joe/2020-04-24",
				Views: 5,
			},
		})
	}()
	<-ds.CallsToInsertPageViews

	router := mux.NewRouter()
	s := defaultServer{
		datastore:              &ds,
		router:                 router,
		csrfMiddleware:         dummyCsrfMiddleware(),
		googleAnalyticsFetcher: &mf,
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/pageViews?path="+url.QueryEscape("/joe/2020-04-24"), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response pageViewResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}
	viewsExpected := 5
	if response.Views != viewsExpected {
		t.Fatalf("unexpected view count: got %v want %v", response.Views, viewsExpected)
	}

	// Wait for an async call to ds.InsertPageViews.
	select {
	case <-ds.CallsToInsertPageViews:
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for call to ds.InsertPageViews")
	}

	var expected = []struct {
		path          string
		viewsExpected int
	}{
		{"/joe/2020-04-17", 8},
		{"/joe/2020-04-24", 11},
		{"/mary/2020-04-17", 25},
	}
	for _, tt := range expected {
		pvr, err := ds.GetPageViews(tt.path)
		if err != nil {
			t.Fatalf("unexpected response from GetPageViews: %v", err)
		}
		if pvr.PageViews != tt.viewsExpected {
			t.Fatalf("unexpected view count for %s - got: %d, want %d", tt.path, pvr.PageViews, tt.viewsExpected)
		}
	}
}
