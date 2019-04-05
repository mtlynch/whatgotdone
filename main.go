package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Page struct {
	Title string
}

type JournalEntry struct {
	Date         string `json:"date"`
	LastModified string `json:"lastModified"`
	Markdown     string `json:"markdown"`
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

var entries = []JournalEntry{
	JournalEntry{
		Date:         "2019-03-22",
		LastModified: "2019-03-22T08:15:22.382Z",
		Markdown:     "Ate some crackers",
	},
	JournalEntry{
		Date:         "2019-03-15",
		LastModified: "2019-03-15T22:06:45.196Z",
		Markdown:     "Took a nap",
	},
	JournalEntry{
		Date:         "2019-03-08",
		LastModified: "2019-03-22T14:59:16.010Z",
		Markdown:     "Watched the movie *The Royal Tenenbaums*.",
	},
}

func entriesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

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
	entries = append(entries, JournalEntry{
		Date:         t.Date,
		LastModified: time.Now().Format(time.RFC3339),
		Markdown:     t.EntryContent,
	})
	resp := SubmitResponse{
		Ok: true,
	}

	log.Printf("Logged a new entry. Current length: %d\n", len(entries))
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

func addTextEntries() {
	a, err := ioutil.ReadFile("2019-03-29.md")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile("2019-04-05.md")
	if err != nil {
		panic(err)
	}
	entries = append(entries,
		JournalEntry{
			Date:         "2019-04-05",
			LastModified: "2019-04-05T21:36:05.333Z",
			Markdown:     string(b),
		})
	entries = append(entries,
		JournalEntry{
			Date:         "2019-03-29",
			LastModified: "2019-03-29T23:34:02.111Z",
			Markdown:     string(a),
		})
}

func main() {
	addTextEntries()
	fs := http.FileServer(http.Dir("client/dist"))
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
