package main

import (
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Title string
	Body  string
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

func main() {
	fs := http.FileServer(http.Dir("client/dist"))
	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.HandleFunc("/", makeHandler(indexHandler))
	log.Fatal(http.ListenAndServe(":3001", nil))
}
