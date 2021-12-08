package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestDraftHandlerWhenUserIsNotLoggedIn(t *testing.T) {
	drafts := []types.JournalEntry{
		{Date: "2019-04-19", LastModified: mustParseTime("2019-04-19T00:00:00Z"), Markdown: "Drove to the zoo"},
	}
	ds := mock.MockDatastore{
		JournalDrafts: drafts,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator:  mockAuthenticator{},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/draft/2019-04-19", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusForbidden {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestDraftHandlerWhenUserTokenIsInvalid(t *testing.T) {
	drafts := []types.JournalEntry{
		{Date: "2019-04-19", LastModified: mustParseTime("2019-04-19T00:00:00Z"), Markdown: "Drove to the zoo"},
	}
	ds := mock.MockDatastore{
		JournalDrafts: drafts,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/draft/2019-04-19", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_invalid", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusForbidden {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestDraftHandlerWhenDateMatches(t *testing.T) {
	drafts := []types.JournalEntry{
		{Date: "2019-04-19", LastModified: mustParseTime("2019-04-19T00:00:00Z"), Markdown: types.EntryContent("Drove to the zoo")},
	}
	ds := mock.MockDatastore{
		JournalDrafts: drafts,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/draft/2019-04-19", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	type response struct {
		Markdown string `json:"markdown"`
	}
	var resp response
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	if resp.Markdown != string(drafts[0].Markdown) {
		t.Fatalf("Unexpected response: got %v want %v", resp.Markdown, drafts[0].Markdown)
	}
}

func TestDraftHandlerReturns404WhenDatastoreReturnsEntryNotFoundError(t *testing.T) {
	entries := []types.JournalEntry{}
	ds := mock.MockDatastore{
		JournalDrafts: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/draft/2019-04-19", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDraftHandlerReturnsBadRequestWhenDateIsInvalid(t *testing.T) {
	entries := []types.JournalEntry{}
	ds := mock.MockDatastore{
		JournalDrafts: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/draft/201904-19", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func mustParseTime(ts string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", ts)
	if err != nil {
		panic(err)
	}
	return t
}

func TestDeleteDraftDeletesMatchingDraft(t *testing.T) {
	ds := mock.MockDatastore{
		JournalDrafts: []types.JournalEntry{
			{
				Author:       "dummyUser",
				Date:         "2019-03-22",
				LastModified: mustParseTime("2019-03-24T00:00:00Z"),
				Markdown:     "Ate some crackers",
			},
			{
				Author:       "dummyUser",
				Date:         "2019-03-15",
				LastModified: mustParseTime("2019-03-15T00:00:00Z"),
				Markdown:     "Took a nap",
			},
			{
				Author:       "dummyUser",
				Date:         "2019-03-08",
				LastModified: mustParseTime("2019-03-09T00:00:00Z"),
				Markdown:     "Watched the movie *The Royal Tenenbaums*.",
			},
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("DELETE", "/api/draft/2019-03-15", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	statusExpected := http.StatusOK
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}

	draftsExpected := []types.JournalEntry{
		{
			Author:       "dummyUser",
			Date:         "2019-03-22",
			LastModified: mustParseTime("2019-03-24T00:00:00Z"),
			Markdown:     "Ate some crackers",
		},
		{
			Author:       "dummyUser",
			Date:         "2019-03-08",
			LastModified: mustParseTime("2019-03-09T00:00:00Z"),
			Markdown:     "Watched the movie *The Royal Tenenbaums*.",
		},
	}
	if !reflect.DeepEqual(ds.JournalDrafts, draftsExpected) {
		t.Fatalf("datastore in wrong state: got %+v want %+v", ds.JournalDrafts, draftsExpected)
	}
}

func TestDeleteDraftReturnsOKForNonExistentEntry(t *testing.T) {
	ds := mock.MockDatastore{
		JournalDrafts: []types.JournalEntry{
			{
				Author:       "dummyUser",
				Date:         "2019-03-22",
				LastModified: mustParseTime("2019-03-24T00:00:00Z"),
				Markdown:     "Ate some crackers",
			},
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("DELETE", "/api/draft/2019-03-15", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	statusExpected := http.StatusOK
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}

	entriesExpected := []types.JournalEntry{
		{
			Author:       "dummyUser",
			Date:         "2019-03-22",
			LastModified: mustParseTime("2019-03-24T00:00:00Z"),
			Markdown:     "Ate some crackers",
		},
	}
	if !reflect.DeepEqual(ds.JournalDrafts, entriesExpected) {
		t.Fatalf("datastore in wrong state: got %+v want %+v", ds.JournalDrafts, entriesExpected)
	}
}

func TestDeleteDraftReturnsBadRequestForInvalidDate(t *testing.T) {
	ds := mock.MockDatastore{
		JournalDrafts: []types.JournalEntry{
			{
				Author:       "dummyUser",
				Date:         "2019-03-22",
				LastModified: mustParseTime("2019-03-24T00:00:00Z"),
				Markdown:     "Ate some crackers",
			},
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	// 2019-03-16 is a Saturday, not a Friday
	req, err := http.NewRequest("DELETE", "/api/draft/2019-03-16", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	statusExpected := http.StatusBadRequest
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}
}
