package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/types"
)

type mockDatastore struct {
	journalEntries []types.JournalEntry
}

func (ds mockDatastore) All() ([]types.JournalEntry, error) {
	return ds.journalEntries, nil
}

func (ds mockDatastore) InsertJournalEntry(types.JournalEntry) error {
	return nil
}

func (ds mockDatastore) Close() error {
	return nil
}

func init() {
	// The handler uses relative paths to find the template file. Switch to the
	// app's root directory so that the relative paths work.
	if err := os.Chdir("../"); err != nil {
		panic(err)
	}
}

func TestEntriesHandler(t *testing.T) {
	entries := []types.JournalEntry{
		types.JournalEntry{Date: "2019-03-22", LastModified: "2019-03-24", Markdown: "Ate some crackers"},
		types.JournalEntry{Date: "2019-03-15", LastModified: "2019-03-15", Markdown: "Took a nap"},
		types.JournalEntry{Date: "2019-03-08", LastModified: "2019-03-09", Markdown: "Watched the movie *The Royal Tenenbaums*."},
	}
	ds := mockDatastore{
		journalEntries: entries,
	}
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

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var response []types.JournalEntry
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Response is not valid JSON: %v", w.Body.String())
	}

	if !reflect.DeepEqual(response, entries) {
		t.Fatalf("Unexpected response: got %v want %v", response, entries)
	}
}
