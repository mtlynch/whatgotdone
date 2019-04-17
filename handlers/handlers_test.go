package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/types"
)

type mockDatastore struct {
}

func (ds mockDatastore) All() ([]types.JournalEntry, error) {
	return []types.JournalEntry{
		types.JournalEntry{Date: "2019-03-22", LastModified: "2019-03-24", Markdown: "Ate some crackers"},
		types.JournalEntry{Date: "2019-03-15", LastModified: "2019-03-15", Markdown: "Took a nap"},
		types.JournalEntry{Date: "2019-03-08", LastModified: "2019-03-09", Markdown: "Watched the movie *The Royal Tenenbaums*."},
	}, nil
}

func (ds mockDatastore) InsertJournalEntry(types.JournalEntry) error {
	return nil
}

func (ds mockDatastore) Close() error {
	return nil
}

func TestEntriesHandler(t *testing.T) {
	ds := mockDatastore{}
	router := mux.NewRouter()
	s := defaultServer{
		datastore: ds,
		router:    router,
	}
	s.routes()

	req, err := http.NewRequest("GET", "/entries", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	s.router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
