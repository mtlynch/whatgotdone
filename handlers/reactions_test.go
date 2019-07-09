package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/types"
)

func TestReactionsGetWhenEntryHasNoReactions(t *testing.T) {
	ds := mockDatastore{
		reactions: []types.Reaction{},
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore: ds,
		router:    router,
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
	reactions := []types.Reaction{
		types.Reaction{Username: "dummyUserA", Symbol: "🎉", Timestamp: "2019-07-09T14:56:29-04:00"},
		types.Reaction{Username: "dummyUserB", Symbol: "👍", Timestamp: "2019-07-09T11:57:02-04:00"},
	}
	ds := mockDatastore{
		reactions: reactions,
	}
	router := mux.NewRouter()
	s := defaultServer{
		datastore: ds,
		router:    router,
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

	if !reflect.DeepEqual(response, reactions) {
		t.Fatalf("Unexpected response: got %v want %v", response, reactions)
	}
}
