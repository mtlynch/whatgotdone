package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestRecentEntriesRejectsInvalidStartAndLimitParameters(t *testing.T) {
	router := mux.NewRouter()
	s := defaultServer{
		entriesReader:  nil,
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
	router := mux.NewRouter()
	s := defaultServer{
		entriesReader:  nil,
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
