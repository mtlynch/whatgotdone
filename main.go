package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title string
}

type JournalEntry struct {
	Id       int    `json:"name"`
	Date     string `json:"date"`
	Markdown string `json:"markdown"`
}

func indexHandler(w http.ResponseWriter, r *http.Request, title string) {
	p := &Page{
		Title: "What Got Done",
	}
	renderTemplate(w, "index.html", p)
}

var templates = template.Must(
	// Use custom delimiters so Go's delimiters don't clash with Vue's.
	template.New("client/dist/index.html").Delims("[[", "]]").ParseFiles(
		"client/dist/index.html"))

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
	entries := []JournalEntry{
		JournalEntry{Id: 0, Date: "2019-03-25", Markdown: "Ate some crackers"},
		JournalEntry{Id: 1, Date: "2019-03-24", Markdown: "Took a nap"},
		JournalEntry{Id: 2, Date: "2019-03-23", Markdown: "Watched a movie!"},
	}

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		panic(err)
	}
}

func main() {
	fs := http.FileServer(http.Dir("client/dist"))
	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.HandleFunc("/entries", entriesHandler)
	http.HandleFunc("/", makeHandler(indexHandler))
	log.Fatal(http.ListenAndServe(":3001", nil))
}
