package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestDraftHandlerWhenUserIsNotLoggedIn(t *testing.T) {
	drafts := []types.JournalEntry{
		{Date: "2019-04-19", LastModified: "2019-04-19", Markdown: "Drove to the zoo"},
	}
	ds := mockDatastore{
		journalDrafts: drafts,
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
		{Date: "2019-04-19", LastModified: "2019-04-19", Markdown: "Drove to the zoo"},
	}
	ds := mockDatastore{
		journalDrafts: drafts,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
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
		{Date: "2019-04-19", LastModified: "2019-04-19", Markdown: "Drove to the zoo"},
	}
	ds := mockDatastore{
		journalDrafts: drafts,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
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

	if resp.Markdown != drafts[0].Markdown {
		t.Fatalf("Unexpected response: got %v want %v", resp.Markdown, drafts[0].Markdown)
	}
}

func TestDraftHandlerReturns404WhenDatastoreReturnsEntryNotFoundError(t *testing.T) {
	entries := []types.JournalEntry{}
	ds := mockDatastore{
		journalDrafts: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
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
	ds := mockDatastore{
		journalDrafts: entries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]string{
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
