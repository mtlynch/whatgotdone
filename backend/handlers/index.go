package handlers

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
)

// indexHandler returns the index.html template and performs some server-side
// rendering of template variables before the Vue frontend renders the page
// client-side.
func (s *defaultServer) indexHandler(pageTitle string) http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		type page struct {
			Title     string
			CsrfToken string
		}
		p := &page{
			CsrfToken: csrf.Token(r),
			Title:     pageTitle,
		}
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}