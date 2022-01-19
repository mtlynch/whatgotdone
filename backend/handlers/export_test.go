package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
	"github.com/mtlynch/whatgotdone/backend/types/export"
)

func TestExportPopulatedUserAccount(t *testing.T) {
	ds := mock.MockDatastore{
		JournalDrafts: []types.JournalEntry{
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-12"),
				LastModified: mustParseTime("2021-11-12T00:00:00Z"),
				Markdown:     "thought about fishing",
			},
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-19"),
				LastModified: mustParseTime("2021-11-19T00:00:00Z"),
				Markdown:     "went to the store today",
			},
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-26"),
				LastModified: mustParseTime("2021-11-20T00:00:00Z"),
				Markdown:     "bought a new car",
			},
		},
		JournalEntries: []types.JournalEntry{
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-12"),
				LastModified: mustParseTime("2021-11-12T00:00:00Z"),
				Markdown:     "thought about fishing",
			},
			{
				Author:       "dummyUserA",
				Date:         types.EntryDate("2021-11-19"),
				LastModified: mustParseTime("2021-11-19T00:00:00Z"),
				Markdown:     "went to the store today",
			},
		},
		Reactions: map[types.Username]map[types.EntryDate][]types.Reaction{
			"dummyUserA": {
				"2021-11-19": []types.Reaction{
					{
						Username:  types.Username("dummyUserB"),
						Symbol:    "👍",
						Timestamp: mustParseTime("2021-11-20T11:57:02Z"),
					},
				},
			},
		},
		UserFollows: map[types.Username][]types.Username{
			types.Username("dummyUserA"): {types.Username("dummyUserC")},
		},
		UserPreferences: map[types.Username]types.Preferences{
			types.Username("dummyUserA"): {
				EntryTemplate: "# My weekly template",
			},
		},
		UserProfile: types.UserProfile{
			AboutMarkdown: "I'm just a dummy user",
			EmailAddress:  "dummy@example.com",
		},
		Usernames: []types.Username{"dummyUserA", "dummyUserB", "dummyUserC"},
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUserA",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/export", nil)
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
	var response export.UserData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	exportExpected := export.UserData{
		Drafts: []export.JournalEntry{
			{
				Date:         types.EntryDate("2021-11-12"),
				LastModified: "2021-11-12T00:00:00Z",
				Markdown:     "thought about fishing",
			},
			{
				Date:         types.EntryDate("2021-11-19"),
				LastModified: "2021-11-19T00:00:00Z",
				Markdown:     "went to the store today",
			},
			{
				Date:         types.EntryDate("2021-11-26"),
				LastModified: "2021-11-20T00:00:00Z",
				Markdown:     "bought a new car",
			},
		},
		Entries: []export.JournalEntry{
			{
				Date:         types.EntryDate("2021-11-12"),
				LastModified: "2021-11-12T00:00:00Z",
				Markdown:     "thought about fishing",
			},
			{
				Date:         types.EntryDate("2021-11-19"),
				LastModified: "2021-11-19T00:00:00Z",
				Markdown:     "went to the store today",
			},
		},
		Reactions: map[types.EntryDate][]export.Reaction{
			"2021-11-19": {
				{
					Username:  types.Username("dummyUserB"),
					Symbol:    "👍",
					Timestamp: "2021-11-20T11:57:02Z",
				},
			},
		},
		Following: []types.Username{types.Username("dummyUserC")},
		Preferences: export.Preferences{
			EntryTemplate: "# My weekly template",
		},
		Profile: export.UserProfile{
			AboutMarkdown: "I'm just a dummy user",
			EmailAddress:  "dummy@example.com",
		},
	}
	if !reflect.DeepEqual(response, exportExpected) {
		t.Fatalf("Unexpected response: got %+v want %+v", response, exportExpected)
	}
}

func TestExportEmptyUserAccount(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUserA",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/export", nil)
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
	var response export.UserData
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	exportExpected := export.UserData{
		Entries:   []export.JournalEntry{},
		Reactions: map[types.EntryDate][]export.Reaction{},
		Drafts:    []export.JournalEntry{},
	}
	if !reflect.DeepEqual(response, exportExpected) {
		t.Fatalf("Unexpected response: got %#v want %#v", response, exportExpected)
	}
}

func TestExportUnauthenticatedAccount(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_A": "dummyUserA",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/export", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock-invalid-token", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)

	if status := w.Code; status != http.StatusForbidden {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}
}
