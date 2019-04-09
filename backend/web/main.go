package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mtlynch/whatgotdone/backend/datastore"
	"github.com/mtlynch/whatgotdone/backend/types"
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
		"../../client/dist/index.html"))

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
		log.Printf("Failed to decode request: %s\n", r.Body)
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
	fs := http.FileServer(http.Dir("../../client/dist"))
	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.HandleFunc("/entries", entriesHandler)
	http.HandleFunc("/api/submit", submitHandler)
	http.HandleFunc("/", makeHandler(indexHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
