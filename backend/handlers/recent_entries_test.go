package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func TestRecentEntriesHandlerSortsByDateThenByModifedTimeInDescendingOrder(t *testing.T) {
	entries := []types.JournalEntry{
		types.JournalEntry{Date: "2019-05-24", LastModified: "2019-05-24T00:00:00.000Z", Markdown: "Rode the bus and saw a movie about ghosts"},
		types.JournalEntry{Date: "2019-05-24", LastModified: "2019-05-23T00:00:00.000Z", Markdown: "Ate some crackers in a bathtub"},
		types.JournalEntry{Date: "2019-05-17", LastModified: "2019-05-17T12:00:00.000Z", Markdown: "Saw a movie about French vanilla"},
		types.JournalEntry{Date: "2019-05-24", LastModified: "2019-05-25T00:00:00.000Z", Markdown: "Read a book about the history of cheese"},
		types.JournalEntry{Date: "2019-05-24", LastModified: "2019-05-25T22:00:00.000Z", Markdown: "Read a pamphlet from The Cat Society"},
		types.JournalEntry{Date: "2019-05-24", LastModified: "2019-05-25T06:00:00.000Z", Markdown: "Read the news today... Oh boy!"},
		types.JournalEntry{Date: "2019-05-17", LastModified: "2019-05-16T00:00:00.000Z", Markdown: "Took a nap and dreamed about chocolate"},
	}
	ds := mockDatastore{
		journalEntries: entries,
		users: []string{
			"bob",
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore: &ds,
		router:    router,
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/recentEntries", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response []recentEntry
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	// For simplicity of the test, all users have username "bob," but in
	// practice, these updates would come from different users.
	expected := []recentEntry{
		recentEntry{Author: "bob", Date: "2019-05-24", Markdown: "Read a pamphlet from The Cat Society"},
		recentEntry{Author: "bob", Date: "2019-05-24", Markdown: "Read the news today... Oh boy!"},
		recentEntry{Author: "bob", Date: "2019-05-24", Markdown: "Read a book about the history of cheese"},
		recentEntry{Author: "bob", Date: "2019-05-24", Markdown: "Rode the bus and saw a movie about ghosts"},
		recentEntry{Author: "bob", Date: "2019-05-24", Markdown: "Ate some crackers in a bathtub"},
		recentEntry{Author: "bob", Date: "2019-05-17", Markdown: "Saw a movie about French vanilla"},
		recentEntry{Author: "bob", Date: "2019-05-17", Markdown: "Took a nap and dreamed about chocolate"},
	}
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}

func TestRecentEntriesHandlerAlwaysPlacesNewDatesAheadOfOldDates(t *testing.T) {
	entries := []types.JournalEntry{
		types.JournalEntry{Date: "2019-09-20", LastModified: "2019-09-20T00:00:00.000Z", Markdown: "Attended an Indie Hackers meetup"},
		types.JournalEntry{Date: "2019-09-13", LastModified: "2019-09-25T00:00:00.000Z", Markdown: "High fived a platypus"},
		types.JournalEntry{Date: "2019-05-17", LastModified: "2019-09-28T12:00:00.000Z", Markdown: "Made a hat out of donuts"},
		types.JournalEntry{Date: "2019-09-06", LastModified: "2019-09-22T00:00:00.000Z", Markdown: "Ate an apple in a single bite"},
	}
	ds := mockDatastore{
		journalEntries: entries,
		users: []string{
			"bob",
		},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore: &ds,
		router:    router,
	}
	s.routes()

	req, err := http.NewRequest("GET", "/api/recentEntries", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response []recentEntry
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Response is not valid JSON: %v", w.Body.String())
	}

	// For simplicity of the test, all users have username "bob," but in
	// practice, these updates would come from different users.
	expected := []recentEntry{
		recentEntry{Author: "bob", Date: "2019-09-20", Markdown: "Attended an Indie Hackers meetup"},
		recentEntry{Author: "bob", Date: "2019-09-13", Markdown: "High fived a platypus"},
		recentEntry{Author: "bob", Date: "2019-09-06", Markdown: "Ate an apple in a single bite"},
		recentEntry{Author: "bob", Date: "2019-05-17", Markdown: "Made a hat out of donuts"},
	}
	if !reflect.DeepEqual(response, expected) {
		t.Fatalf("Unexpected response: got %v want %v", response, expected)
	}
}
