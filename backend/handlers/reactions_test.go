package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/datastore/mock"
	"github.com/mtlynch/whatgotdone/backend/types"
)

// Create a dummy CSRF middleware that never rejects HTTP requests.
func dummyCsrfMiddleware() httpMiddlewareHandler {
	return func(h http.Handler) http.Handler {
		return h
	}
}

func TestReactionsGetWhenEntryHasNoReactions(t *testing.T) {
	ds := mock.MockDatastore{
		Reactions: map[types.Username]map[types.EntryDate][]types.Reaction{},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/reactions/entry/dummyUser/2019-07-12", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response []types.Reaction
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	if len(response) != 0 {
		t.Fatalf("Unexpected response: got %v want []", response)
	}
}

func TestReactionsGetWhenEntryHasTwoReactions(t *testing.T) {
	ds := mock.MockDatastore{
		Reactions: map[types.Username]map[types.EntryDate][]types.Reaction{
			"dummyUser": {
				"2019-07-12": {
					{Username: "dummyUserA", Symbol: "🎉", Timestamp: "2019-07-09T14:56:29-04:00"},
					{Username: "dummyUserB", Symbol: "👍", Timestamp: "2019-07-09T11:57:02-04:00"},
				},
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

	req, err := http.NewRequest("GET", "/api/reactions/entry/dummyUser/2019-07-12", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response []types.Reaction
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	expected := []types.Reaction{
		{Username: "dummyUserA", Symbol: "🎉", Timestamp: "2019-07-09T14:56:29-04:00"},
		{Username: "dummyUserB", Symbol: "👍", Timestamp: "2019-07-09T11:57:02-04:00"},
	}
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}

func TestReactionsGetWhenEntryAuthorIsUndefined(t *testing.T) {
	ds := mock.MockDatastore{
		Reactions: map[types.Username]map[types.EntryDate][]types.Reaction{},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/reactions/entry/undefined/2019-07-12", nil)
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

func TestReactionsPostStoresValidReaction(t *testing.T) {
	ds := mock.MockDatastore{
		Reactions: map[types.Username]map[types.EntryDate][]types.Reaction{},
	}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "reactionSymbol": "👍" }`)
	req, err := http.NewRequest("POST", "/api/reactions/entry/dummyUserA/2019-04-19", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	reactions, ok := ds.Reactions["dummyUserA"]["2019-04-19"]
	if !ok {
		t.Fatalf("no reactions list found at dummyUserA/2019-04-19 in mock datastore: %+v", ds.Reactions)
	}
	if len(reactions) != 1 {
		t.Fatalf("unexpected reaction count: got %v (%v) want %v",
			len(reactions), reactions, 1)
	}
	if reactions[0].Username != "dummyUserC" {
		t.Fatalf("unexpected username in reaction: got %v want %v",
			reactions[0].Username, "dummyUserC")
	}
	if reactions[0].Symbol != "👍" {
		t.Fatalf("unexpected symbol in reaction: got [%v] want [%v]",
			reactions[0].Symbol, "👍")
	}
}

func TestReactionsPostRejectsRequestWithMissingSymbolField(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "dummyNonExistentFieldName": "👍" }`)
	req, err := http.NewRequest("POST", "/api/reactions/entry/dummyUserA/2019-04-19", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestReactionsRejectsInvalidReactionSymbol(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "reactionSymbol": "!" }`)
	req, err := http.NewRequest("POST", "/api/reactions/entry/dummyUserA/2019-04-19", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestReactionsPostRejectsRequestWhenUsernameIsUndefined(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator: mockAuthenticator{
			tokensToUsers: map[string]types.Username{
				"mock_token_C": "dummyUserC",
			},
		},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "reactionSymbol": "👍" }`)
	req, err := http.NewRequest("POST", "/api/reactions/entry/undefined/2019-04-19", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", fmt.Sprintf("%s=mock_token_C", userKitAuthCookieName))

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestReactionsPostRejectsRequestWhenUserIsNotLoggedIn(t *testing.T) {
	ds := mock.MockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		authenticator:  mockAuthenticator{},
		datastore:      &ds,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	requestBody := []byte(`{ "symbol": "👍" }`)
	req, err := http.NewRequest("POST", "/api/reactions/entry/dummyUser/2019-04-19", bytes.NewBuffer(requestBody))
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
