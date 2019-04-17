package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mtlynch/whatgotdone/datastore"
	"github.com/mtlynch/whatgotdone/types"
)

type (
	page struct {
		Title string
	}

	submitRequest struct {
		Date         string `json:"date"`
		EntryContent string `json:"entryContent"`
	}

	submitResponse struct {
		Ok bool `json:"ok"`
	}
)

var templates = template.Must(
	// Use custom delimiters so Go's delimiters don't clash with Vue's.
	template.New("index.html").Delims("[[", "]]").ParseFiles(
		"./web/frontend/dist/index.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	enableCsp(&w)
	p := &page{
		Title: "What Got Done",
	}
	renderTemplate(w, "index.html", p)
}

func EntriesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	ds := datastore.New()
	defer ds.Close()
	entries, err := ds.All()
	if err != nil {
		log.Printf("Failed to retrieve entries: %s", err)
		return
	}

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		panic(err)
	}
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var t submitRequest
	err := decoder.Decode(&t)
	if err != nil {
		log.Printf("Failed to decode request: %s", r.Body)
	}
	j := types.JournalEntry{
		Date:         t.Date,
		LastModified: time.Now().Format(time.RFC3339),
		Markdown:     t.EntryContent,
	}
	ds := datastore.New()
	defer ds.Close()
	err = ds.InsertJournalEntry(j)
	if err != nil {
		log.Printf("Failed to insert journal entry: %s", err)
	}
	resp := submitResponse{
		Ok: true,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *page) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func enableCsp(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Security-Policy", "default-src 'self'")
}
