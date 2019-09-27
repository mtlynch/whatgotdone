package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/csrf"
)

type page struct {
	Title     string
	CsrfToken string
}

var (
	fs = http.Dir("./frontend/dist")
	// Use custom delimiters so Go's delimiters don't clash with Vue's.
	indexTemplate = template.Must(template.New("index.html").Delims("[[", "]]").
			ParseFiles("./frontend/dist/index.html"))
)

// serveStaticPage serves any static file under `./frontend/dist` or if said
// file does not exist then it returns the index.html template and performs some
// server-side rendering of template variables before the Vue frontend renders
// the page client-side.
func (s defaultServer) serveStaticPage() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Open the file
		file, err := fs.Open(r.URL.Path)
		if err == nil {
			// Gather the information about the file
			stats, err := file.Stat()
			if err != nil {
				log.Printf("Failed to retrieve the information of %s from the file system: %s", r.URL.Path, err)
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			// We can't serve directories
			if !stats.IsDir() {
				// Copy the file back to the user
				if _, err = io.Copy(rw, file); err != nil {
					log.Printf("Failed to copy the file %s to the response writer: %s", r.URL.Path, err)
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		}
		// If our issue was anything other than the file not existing
		if err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to retrieve the file %s from the file system: %s", r.URL.Path, err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		serveIndexPage(rw, r)
	}
}

// serveIndexPage returns the file `./frontend/dist/index.html` rendered by the
// golang templating engine.
func serveIndexPage(rw http.ResponseWriter, r *http.Request) {
	if err := indexTemplate.ExecuteTemplate(rw, "index.html", page{
		CsrfToken: csrf.Token(r),
		Title:     getPageTitle(r),
	}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// getPageTitle returns the <title> value of the page. By default it's
// "What Got Done" but if the date or username are available, we prepend those
// to the title, so that it can be "username's What Got Done for the week of YYYY-MM-DD".
func getPageTitle(r *http.Request) string {
	t := "What Got Done"

	entryAuthor, err := usernameFromRequestPath(r)
	if err == nil {
		t = entryAuthor + "'s " + t
	}

	date, err := dateFromRequestPath(r)
	if err == nil {
		t += " for the week of " + date
	}

	return t
}
