package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/datastore"
	"github.com/mtlynch/whatgotdone/types"
)

type Page struct {
	Title string
}

func indexHandler(w http.ResponseWriter, r *http.Request, title string) {
	enableCsp(&w)
	p := &Page{
		Title: "What Got Done",
	}
	renderTemplate(w, "index.html", p)
}

var templates = template.Must(
	// Use custom delimiters so Go's delimiters don't clash with Vue's.
	template.New("index.html").Delims("[[", "]]").ParseFiles(
		"./web/frontend/dist/index.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, r.URL.Path)
	}
}

func entriesHandler(w http.ResponseWriter, r *http.Request) {
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

type SubmitRequest struct {
	Date         string `json:"date"`
	EntryContent string `json:"entryContent"`
}

type SubmitResponse struct {
	Ok bool `json:"ok"`
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var t SubmitRequest
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
	resp := SubmitResponse{
		Ok: true,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func enableCsp(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Security-Policy", "default-src 'self'")
}

func main() {
	fs := http.FileServer(http.Dir("./web/frontend/dist"))
	r := mux.NewRouter()
	r.PathPrefix("/js").Handler(fs)
	r.PathPrefix("/css").Handler(fs)
	r.HandleFunc("/entries", entriesHandler)
	r.HandleFunc("/api/submit", submitHandler)
	r.PathPrefix("/").HandlerFunc(makeHandler(indexHandler))
	http.Handle("/", r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("Listening on %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
