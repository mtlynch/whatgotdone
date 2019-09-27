package handlers

import (
	"net/http"
	"text/template"

	"github.com/gorilla/csrf"
)

// indexHandler returns the index.html template and performs some server-side
// rendering of template variables before the Vue frontend renders the page
// client-side.
func (s *defaultServer) indexHandler() http.HandlerFunc {
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
			Title:     getPageTitle(r),
		}
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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