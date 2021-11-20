package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"testing"

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
		PageViewCounts: []ga.PageViewCount{
			{
				Path:  "/jimmy123/2020-01-17",
				Views: 5,
			},
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

type mockGoogleAnalyticsFetcher struct {
	PageViewCounts []ga.PageViewCount
}

func (f mockGoogleAnalyticsFetcher) PageViewsByPath(_, _ string) ([]ga.PageViewCount, error) {
	return f.PageViewCounts, nil
}

func TestRefreshGoogleAnalytics(t *testing.T) {
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
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:              &ds,
		router:                 router,
		csrfMiddleware:         dummyCsrfMiddleware(),
		googleAnalyticsFetcher: mf,
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/tasks/refreshGoogleAnalytics", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Appengine-Cron", "true")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := []ga.PageViewCount{
		{
			Path:  "/joe/2020-04-17",
			Views: 8,
		},
		{
			Path:  "/joe/2020-04-24",
			Views: 11,
		},
		{
			Path:  "/mary/2020-04-17",
			Views: 25,
		},
	}
	sort.Slice(ds.PageViewCounts, func(i, j int) bool { return ds.PageViewCounts[i].Path < ds.PageViewCounts[j].Path })
	if !reflect.DeepEqual(ds.PageViewCounts, expected) {
		t.Fatalf("Unexpected response: got %v want %v", ds.PageViewCounts, expected)
	}
}

func TestRefreshGoogleAnalyticsRejectsExternalRequests(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:              &ds,
		router:                 router,
		csrfMiddleware:         dummyCsrfMiddleware(),
		googleAnalyticsFetcher: mockGoogleAnalyticsFetcher{},
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/tasks/refreshGoogleAnalytics", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Omit AppEngine header

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusForbidden {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}
