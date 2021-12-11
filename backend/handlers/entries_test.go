package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestEntriesHandler(t *testing.T) {
	entries := []types.JournalEntry{
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
	}
	ds := mock.MockDatastore{
		JournalEntries: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/entries/dummyUser", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response []types.JournalEntry
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	if !reflect.DeepEqual(response, entries) {
		t.Fatalf("Unexpected response: got %v want %v", response, entries)
	}
}
func TestEntriesHandlerWhenUserHasNoEntries(t *testing.T) {
	ds := mock.MockDatastore{
		JournalEntries: []types.JournalEntry{},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/entries/dummyUser", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response []types.JournalEntry
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	var expected []types.JournalEntry
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}

func TestEntriesHandlerReturnsBadRequestWhenUsernameIsBlank(t *testing.T) {
	entries := []types.JournalEntry{}
	ds := mock.MockDatastore{
		JournalEntries: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/entries", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestEntriesHandlerReturnsNotFoundWhenUsernameHasNoEntries(t *testing.T) {
	entries := []types.JournalEntry{}
	ds := mock.MockDatastore{
		JournalEntries: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/entries", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPutEntryRejectsEmptyEntry(t *testing.T) {
	ds := mock.MockDatastore{}
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

	req, err := http.NewRequest(
		"PUT",
		"/api/entry/2019-03-15",
		strings.NewReader(`{"entryContent": ""}`))
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

func TestDeleteEntryDeletesMatchingEntry(t *testing.T) {
	ds := mock.MockDatastore{
		JournalEntries: []types.JournalEntry{
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

	req, err := http.NewRequest("DELETE", "/api/entry/2019-03-15", nil)
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
		{
			Author:       "dummyUser",
			Date:         "2019-03-08",
			LastModified: mustParseTime("2019-03-09T00:00:00Z"),
			Markdown:     "Watched the movie *The Royal Tenenbaums*.",
		},
	}
	if !reflect.DeepEqual(ds.JournalEntries, entriesExpected) {
		t.Fatalf("datastore in wrong state: got %+v want %+v", ds.JournalEntries, entriesExpected)
	}
}

func TestDeleteEntryReturnsOKForNonExistentEntry(t *testing.T) {
	ds := mock.MockDatastore{
		JournalEntries: []types.JournalEntry{
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

	req, err := http.NewRequest("DELETE", "/api/entry/2019-03-15", nil)
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
	if !reflect.DeepEqual(ds.JournalEntries, entriesExpected) {
		t.Fatalf("datastore in wrong state: got %+v want %+v", ds.JournalEntries, entriesExpected)
	}
}

func TestDeleteEntryReturnsBadRequestForInvalidDate(t *testing.T) {
	ds := mock.MockDatastore{
		JournalEntries: []types.JournalEntry{
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
	req, err := http.NewRequest("DELETE", "/api/entry/2019-03-16", nil)
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
