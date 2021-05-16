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

type mockEntriesReader struct {
	entries []entries.RecentEntry
}

func (mer mockEntriesReader) Recent(start, limit int) ([]entries.RecentEntry, error) {
	return mer.entries, nil
}

func (mer mockEntriesReader) RecentFollowing(username types.Username, start, limit int) ([]entries.RecentEntry, error) {
	return mer.entries, nil
}

func TestRecentEntriesHandlerReturnsRecentEntries(t *testing.T) {
	recentEntries := []entries.RecentEntry{
		{Author: "alan", Date: "2019-05-24", Markdown: "Read a pamphlet from The Cat Society"},
		{Author: "janie", Date: "2019-05-24", Markdown: "Read the news today... Oh boy!"},
		{Author: "carla", Date: "2019-05-24", Markdown: "Read a book about the history of cheese"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Author: "ted", Date: "2019-05-24", Markdown: "Ate some crackers in a bathtub"},
		{Author: "joe", Date: "2019-05-17", Markdown: "Saw a movie about French vanilla"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Took a nap and dreamed about chocolate"},
	}
	mer := mockEntriesReader{
		entries: recentEntries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &mockDatastore{},
		entriesReader:  &mer,
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

	var response []entryPublic
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	// For simplicity of the test, all users have username "bob," but in
	// practice, these updates would come from different users.
	expected := []entryPublic{
		{Author: "alan", Date: "2019-05-24", Markdown: "Read a pamphlet from The Cat Society"},
		{Author: "janie", Date: "2019-05-24", Markdown: "Read the news today... Oh boy!"},
		{Author: "carla", Date: "2019-05-24", Markdown: "Read a book about the history of cheese"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Author: "ted", Date: "2019-05-24", Markdown: "Ate some crackers in a bathtub"},
		{Author: "joe", Date: "2019-05-17", Markdown: "Saw a movie about French vanilla"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Took a nap and dreamed about chocolate"},
	}
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}

func TestRecentEntriesRejectsInvalidStartAndLimitParameters(t *testing.T) {
	recentEntries := []entries.RecentEntry{
		{Author: "bob", Date: "2019-05-24", Markdown: "Read a pamphlet from The Cat Society"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Read the news today... Oh boy!"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Read a book about the history of cheese"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Ate some crackers in a bathtub"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Saw a movie about French vanilla"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Took a nap and dreamed about chocolate"},
	}
	mer := mockEntriesReader{
		entries: recentEntries,
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &mockDatastore{},
		entriesReader:  &mer,
		router:         router,
		csrfMiddleware: dummyCsrfMiddleware(),
	}

	s.routes()
	var tests = []struct {
		explanation string
		start       string
		limit       string
	}{
		{
			"rejects invalid start",
			"invalid-start-value",
			"3",
		},
		{
			"rejects negative start",
			"-5",
			"3",
		},
		{
			"rejects invalid limit value",
			"2",
			"invalid-limit-value",
		},
		{
			"rejects negative limit",
			"2",
			"-10",
		},
		{
			"rejects zero limit",
			"2",
			"0",
		},
	}
	for _, tt := range tests {
		req, err := http.NewRequest("GET", fmt.Sprintf("/api/recentEntries?start=%s&limit=%s", tt.start, tt.limit), nil)
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
}

func TestRecentEntriesHandlerReturnsEmptyArrayWhenEntriesReaderIsEmpty(t *testing.T) {
	mer := mockEntriesReader{
		entries: []entries.RecentEntry{},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore:      &mockDatastore{},
		entriesReader:  &mer,
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
