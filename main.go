package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

var realEntry = `
# [IsItKeto](https://isitketo.org)

* Migrated to Circle
* Migrated from Selenium to Cypress
* Fixed generation of structured data
* Migrated admin interface to Circle
* Wrote a script to find missing crosslinks
* Wrote a script to download the datastore locally
* Commissioned cartoonist to create new keto meter graphics
* Tried reaching out to KetoSizeMe about visualizations but no dice
* Checked in script for making Tweet spreadsheet
* Fixed a bunch of errors in the NYT data
* At 439 Twitter followers (forgot to check last week)
* Still at 173 foods

# Conferences

* Mostly finished Good Developers Bad Tests presentation

# [Zestful](https://zestfuldata.com)

* Researched Sizzle app as potential customer
* Reached out to author of forked NYT repo
* Reached out to UPenn student and set up private Zestful instance
  * Added README instructions for creating a new private instance
* Attempted to retrain, but accuracy dropped
  * Still, 200+ new labeled examples
* Stripped out “homemade or store-bought” from input

# What Got Done App

* Got a basic version working
  * Successfully runs Golang on the backend and Vue.js on the frontend and they can communicate
  * No database or deployments yet

# [mtlynch.io](https://mtlynch.io)

* Removed G+ social sharing button

# [Dusty VCR](https://dustyvcr.com)

* Moved RSS feed to Feedburner

# Misc

* Finished Indie Jewel Thieves podcast
* Talked to reporter
* Finished working with accountant on taxes
* Attended first beekeeping workshop
* Submitted pull requests to Cypress
  * [Load test spec at runtime rather than Docker image build time](https://github.com/cypress-io/cypress-example-docker-compose/pull/5)
  * [Adding Cypress artifacts folders to .gitignore](https://github.com/cypress-io/cypress-example-docker-compose/pull/4)
`

func entriesHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	entries := []JournalEntry{
		JournalEntry{
			Date:         "2019-03-29",
			LastModified: "2019-03-29T23:34:02.111Z",
			Markdown:     realEntry,
		},
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

	if err := json.NewEncoder(w).Encode(entries); err != nil {
		panic(err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func enableCsp(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Security-Policy", "default-src 'self'")
}

func main() {
	fs := http.FileServer(http.Dir("client/dist"))
	http.Handle("/css/", fs)
	http.Handle("/js/", fs)
	http.HandleFunc("/entries", entriesHandler)
	http.HandleFunc("/", makeHandler(indexHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
