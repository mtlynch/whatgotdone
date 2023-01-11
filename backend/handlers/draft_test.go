package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/datastore/sqlite"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestDraftHandlerWhenUserIsNotLoggedIn(t *testing.T) {
	ds := sqlite.New(":memory:")
	ds.InsertDraft("dummyuser", types.JournalEntry{
		Date:         "2019-04-19",
		LastModified: mustParseTime("2019-04-19T00:00:00Z"),
		Markdown:     "Drove to the zoo",
	})

	router := mux.NewRouter()
	s := defaultServer{
		authenticator:  mockAuthenticator{},
		datastore:      ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/draft/2019-04-19", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusUnauthorized {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestDraftHandlerWhenUserTokenIsInvalid(t *testing.T) {
	ds := sqlite.New(":memory:")
	ds.InsertDraft("dummyuser", types.JournalEntry{
		Date:         "2019-04-19",
		LastModified: mustParseTime("2019-04-19T00:00:00Z"),
		Markdown:     "Drove to the zoo",
	})
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusUnauthorized {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}

func TestDraftHandlerWhenDateMatches(t *testing.T) {
	ds := sqlite.New(":memory:")
	ds.InsertDraft("dummyUser", types.JournalEntry{
		Date:         "2019-04-19",
		LastModified: mustParseTime("2019-04-19T00:00:00Z"),
		Markdown:     "Drove to the zoo",
	})
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

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

	expected := "Drove to the zoo"
	if resp.Markdown != expected {
		t.Fatalf("Unexpected response: got %v want %v", resp.Markdown, expected)
	}
}

func TestDraftHandlerReturns404WhenDatastoreReturnsEntryNotFoundError(t *testing.T) {
	ds := sqlite.New(":memory:")
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDraftHandlerReturnsBadRequestWhenDateIsInvalid(t *testing.T) {
	ds := sqlite.New(":memory:")
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

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

func TestPutDraftRejectsEmptyDraft(t *testing.T) {
	ds := sqlite.New(":memory:")
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest(
		"PUT",
		"/api/draft/2019-03-15",
		strings.NewReader(`{"entryContent": ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_A", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	statusExpected := http.StatusBadRequest
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}
}

func TestDeleteDraftDeletesMatchingDraft(t *testing.T) {
	ds := sqlite.New(":memory:")
	ds.InsertDraft("dummyUser", types.JournalEntry{
		Author:       "dummyUser",
		Date:         "2019-03-22",
		LastModified: mustParseTime("2019-03-24T00:00:00Z"),
		Markdown:     "Ate some crackers",
	})
	ds.InsertDraft("dummyUser", types.JournalEntry{
		Author:       "dummyUser",
		Date:         "2019-03-15",
		LastModified: mustParseTime("2019-03-15T00:00:00Z"),
		Markdown:     "Took a nap",
	})
	ds.InsertDraft("dummyUser", types.JournalEntry{
		Author:       "dummyUser",
		Date:         "2019-03-08",
		LastModified: mustParseTime("2019-03-09T00:00:00Z"),
		Markdown:     "Watched the movie *The Royal Tenenbaums*.",
	})

	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

	statusExpected := http.StatusOK
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}

	// Verify deleted draft is gone.
	_, err = ds.GetDraft("dummyUser", types.EntryDate("2019-03-15"))
	if _, ok := err.(datastore.DraftNotFoundError); !ok {
		t.Fatalf("expected entry to be missing, but got: %v", err)
	}

	// Verify other drafts are still there.
	_, err = ds.GetDraft("dummyUser", types.EntryDate("2019-03-22"))
	if err != nil {
		t.Fatalf("unexpected error retrieving draft: %v", err)
	}
	_, err = ds.GetDraft("dummyUser", types.EntryDate("2019-03-08"))
	if err != nil {
		t.Fatalf("unexpected error retrieving draft: %v", err)
	}
}

func TestDeleteDraftReturnsOKForNonExistentEntry(t *testing.T) {
	ds := sqlite.New(":memory:")
	ds.InsertDraft("dummyUser", types.JournalEntry{
		Author:       "dummyUser",
		Date:         "2019-03-22",
		LastModified: mustParseTime("2019-03-24T00:00:00Z"),
		Markdown:     "Ate some crackers",
	})
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

	statusExpected := http.StatusOK
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}

	// Verify draft is still there.
	_, err = ds.GetDraft("dummyUser", types.EntryDate("2019-03-22"))
	if err != nil {
		t.Fatalf("unexpected error retrieving draft: %v", err)
	}
}

func TestDeleteDraftReturnsBadRequestForInvalidDate(t *testing.T) {
	ds := sqlite.New(":memory:")
	ds.InsertDraft("dummyUser", types.JournalEntry{
		Author:       "dummyUser",
		Date:         "2019-03-22",
		LastModified: mustParseTime("2019-03-24T00:00:00Z"),
		Markdown:     "Ate some crackers",
	})
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUser",
			},
		},
		datastore:      ds,
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
	s.Router().ServeHTTP(w, req)

	statusExpected := http.StatusBadRequest
	if status := w.Code; status != statusExpected {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, statusExpected)
	}
}
