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
	"github.com/mtlynch/whatgotdone/backend/entries"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestRecentEntriesRejectsInvalidStartAndLimitParameters(t *testing.T) {
	journalEntries := []types.JournalEntry{
		{Date: "2019-05-10", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
	}
	ds := mockDatastore{
		journalEntries: journalEntries,
		users: []string{
			"bob",
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		entriesReader:  entries.NewReader(&ds),
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()
	var tests = []struct {
		explanation     string
		start           string
		limit           string
		statusExpected  int
		entriesExpected []entryPublic
	}{
		{
			"rejects invalid start",
			"invalid-start-value",
			"3",
			http.StatusBadRequest,
			[]entryPublic{},
		},
		{
			"rejects negative start",
			"-5",
			"3",
			http.StatusBadRequest,
			[]entryPublic{},
		},
		{
			"rejects invalid limit value",
			"2",
			"invalid-limit-value",
			http.StatusBadRequest,
			[]entryPublic{},
		},
		{
			"rejects negative limit",
			"2",
			"-10",
			http.StatusBadRequest,
			[]entryPublic{},
		},
		{
			"rejects zero limit",
			"2",
			"0",
			http.StatusBadRequest,
			[]entryPublic{},
		},
	}
	for _, tt := range tests {
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/recentEntries?start=%s&limit=%s", tt.start, tt.limit), nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		s.router.ServeHTTP(w, req)

		if status := w.Code; status != tt.statusExpected {
			t.Fatalf("handler returned wrong status code: got %v want %v",
				status, tt.statusExpected)
		}
		if tt.statusExpected != http.StatusOK {
			continue
		}

		var response []entryPublic
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Response is not valid JSON: %v", w.Body.String())
		}

		if !reflect.DeepEqual(response, tt.entriesExpected) {
			t.Fatalf("%s: Unexpected response: got %v want %v", tt.explanation, response, tt.entriesExpected)
		}
	}
}

func TestRecentEntriesHandlerReturnsEmptyArrayWhenDatastoreIsEmpty(t *testing.T) {
	ds := mockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &ds,
		entriesReader:  entries.NewReader(&ds),
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/recentEntries?start=0&limit=15", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	response := strings.TrimSpace(w.Body.String())
	want := "[]"
	if response != want {
		t.Fatalf("Unexpected response: got %v want %v", response, want)
	}
}
