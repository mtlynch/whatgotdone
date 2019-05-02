package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/datastore"
	"github.com/mtlynch/whatgotdone/types"
)

func (s *defaultServer) indexHandler() http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./web/frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		enableCsp(&w)

		type page struct {
			Title string
		}
		p := &page{
			Title: "What Got Done",
		}
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *defaultServer) entriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		entries, err := s.datastore.All(usernameFromRequestPath(r))
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			return
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) entryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		date, err := dateFromRequestPath(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		j, err := s.datastore.Get(usernameFromRequestPath(r), date)
		if err != nil {
			if _, ok := err.(datastore.EntryNotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Printf("Failed to retrieve entry: %s", err)
			http.Error(w, "Failed to retrieve entry", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(j); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) submitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}
		decoder := json.NewDecoder(r.Body)

		type submitRequest struct {
			Date         string `json:"date"`
			EntryContent string `json:"entryContent"`
		}

		type submitResponse struct {
			Ok bool `json:"ok"`
		}

		var t submitRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", r.Body)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
		}
		j := types.JournalEntry{
			Date:         t.Date,
			LastModified: time.Now().Format(time.RFC3339),
			Markdown:     t.EntryContent,
		}
		// TODO: Replace based on the logged-in user.
		const username string = "michael"
		err = s.datastore.Insert(username, j)
		if err != nil {
			log.Printf("Failed to insert journal entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
		}
		resp := submitResponse{
			Ok: true,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

// Catchall for when no API route matches.
func (s *defaultServer) apiRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func enableCsp(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Security-Policy", "default-src 'self'")
}

func usernameFromRequestPath(r *http.Request) string {
	return mux.Vars(r)["username"]
}

func dateFromRequestPath(r *http.Request) (string, error) {
	date := mux.Vars(r)["date"]
	dateFormat := "2006-01-02"
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", errors.New("Invalid date format: must be YYYY-MM-DD")
	}
	return date, nil
}
