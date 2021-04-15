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

func TestRecentEntriesHandlerSortsByDateThenByModifedTimeInDescendingOrder(t *testing.T) {
	journalEntries := []types.JournalEntry{
		{Date: "2019-05-24", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Date: "2019-05-24", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
		{Date: "2019-05-17", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
		{Date: "2019-05-24", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
		{Date: "2019-05-24", LastModified: "2019-05-25T22:00:00.000Z", Markdown: "Read a pamphlet from The Cat Society"},
		{Date: "2019-05-24", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
		{Date: "2019-05-17", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
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
		{Author: "bob", Date: "2019-05-24", Markdown: "Read a pamphlet from The Cat Society"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Read the news today... Oh boy!"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Read a book about the history of cheese"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Rode the bus and saw a movie about ghosts"},
		{Author: "bob", Date: "2019-05-24", Markdown: "Ate some crackers in a bathtub"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Saw a movie about French vanilla"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Took a nap and dreamed about chocolate"},
	}
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}

func TestRecentEntriesHandlerAlwaysPlacesNewDatesAheadOfOldDates(t *testing.T) {
	journalEntries := []types.JournalEntry{
		{Date: "2019-05-17", LastModified: "2019-09-28T12:00:00.000Z", Markdown: "Made a hat out of donuts from the cloud in the sky"},
		{Date: "2019-09-20", LastModified: "2019-09-25T00:00:00.000Z", Markdown: "High fived a platypus when the apple hits the pie."},
		{Date: "2019-09-06", LastModified: "2019-09-22T00:00:00.000Z", Markdown: "Ate an apple in a single bite of choclate"},
		{Date: "2019-09-20", LastModified: "2019-09-20T00:00:00.000Z", Markdown: "Attended an Indie Hackers meetup"},
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
		{Author: "bob", Date: "2019-09-20", Markdown: "High fived a platypus when the apple hits the pie."},
		{Author: "bob", Date: "2019-09-20", Markdown: "Attended an Indie Hackers meetup"},
		{Author: "bob", Date: "2019-09-06", Markdown: "Ate an apple in a single bite of choclate"},
		{Author: "bob", Date: "2019-05-17", Markdown: "Made a hat out of donuts from the cloud in the sky"},
	}
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}

func TestRecentEntriesObservesStartAndLimitParameters(t *testing.T) {
	journalEntries := []types.JournalEntry{
		{Date: "2019-05-10", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
		{Date: "2019-05-03", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
		{Date: "2019-04-26", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
		{Date: "2019-04-19", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
		{Date: "2019-04-12", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
		{Date: "2019-04-05", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
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
			"observes valid start and limit values",
			"1",
			"3",
			http.StatusOK,
			[]entryPublic{
				{Author: "bob", Date: "2019-05-03", Markdown: "Took a nap and dreamed about chocolate"},
				{Author: "bob", Date: "2019-04-26", Markdown: "Read a book about the history of cheese"},
				{Author: "bob", Date: "2019-04-19", Markdown: "Saw a movie about French vanilla"},
			},
		},
		{
			"accepts large ranges",
			"0",
			"500",
			http.StatusOK,
			[]entryPublic{
				{Author: "bob", Date: "2019-05-10", Markdown: "Read the news today... Oh boy!"},
				{Author: "bob", Date: "2019-05-03", Markdown: "Took a nap and dreamed about chocolate"},
				{Author: "bob", Date: "2019-04-26", Markdown: "Read a book about the history of cheese"},
				{Author: "bob", Date: "2019-04-19", Markdown: "Saw a movie about French vanilla"},
				{Author: "bob", Date: "2019-04-12", Markdown: "Ate some crackers in a bathtub"},
				{Author: "bob", Date: "2019-04-05", Markdown: "Rode the bus and saw a movie about ghosts"},
			},
		},
		{
			"returns empty for start beyond size of total response",
			"500",
			"5",
			http.StatusOK,
			[]entryPublic{},
		},
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
