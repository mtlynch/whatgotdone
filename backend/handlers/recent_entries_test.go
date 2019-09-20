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
