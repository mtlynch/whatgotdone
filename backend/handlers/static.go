package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/gorilla/csrf"
	"github.com/mtlynch/whatgotdone/backend/datastore"
)

const frontendRootDir = "./frontend/dist"
const frontendIndexFilename = "index.html"

// serveStaticResource serves any static file under `./frontend/dist` or if said
// file does not exist then it returns the index.html template and performs some
// server-side rendering of template variables before the Vue frontend renders
// the page client-side.
func (s defaultServer) serveStaticResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fs := http.Dir(frontendRootDir)
		file, err := fs.Open(r.URL.Path)
		if os.IsNotExist(err) {
			log.Printf("%s does not exist on the file system: %s", r.URL.Path, err)
			http.Error(w, "Failed to find file: "+r.URL.Path, http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("Failed to retrieve the file %s from the file system: %s", r.URL.Path, err)
			http.Error(w, "Failed to find file: "+r.URL.Path, http.StatusNotFound)
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Printf("Failed to retrieve the information of %s from the file system: %s", r.URL.Path, err)
			http.Error(w, "Failed to serve: "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		if stat.IsDir() {
			log.Printf("%s is a directory", r.URL.Path)
			http.Error(w, "Failed to find file: "+r.URL.Path, http.StatusNotFound)
			return
		}

		http.ServeFile(w, r, path.Join(frontendRootDir, r.URL.Path))
	}
}

// serveEntyOr404 tries to find an entry associated with the given route or
// returns a 404 if there's no associated entry.
func (s defaultServer) serveEntryOr404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			serve404(w, r)
			return
		}
		if exists, err := s.userExists(username); err != nil || !exists {
			serve404(w, r)
			return
		}

		date, err := dateFromRequestPath(r)
		if err != nil {
			serve404(w, r)
			return
		}
		_, err = s.datastore.GetEntry(username, date)
		if _, ok := err.(datastore.EntryNotFoundError); ok {
			serve404(w, r)
			return
		}
		serveIndexPage(w, r)
	}
}

// serveUserProfileOr404 tries to find a valid user associated with the given
// route or returns a 404 if there's no associated user.
func (s defaultServer) serveUserProfileOr404() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := usernameFromRequestPath(r)
		if err != nil {
			serve404(w, r)
			return
		}
		if exists, err := s.userExists(username); err != nil || !exists {
			serve404(w, r)
			return
		}
		serveIndexPage(w, r)
	}
}

// serveIndexPage returns the file `./frontend/dist/index.html` rendered by the
// golang templating engine.
func serveIndexPage(w http.ResponseWriter, r *http.Request) {
	type page struct {
		Title         string
		Description   string
		CsrfToken     string
		OpenGraphType string
	}
	// Use custom delimiters so Go's delimiters don't clash with Vue's.
	indexTemplate := template.Must(template.New(frontendIndexFilename).Delims("[[", "]]").
		ParseFiles(path.Join(frontendRootDir, frontendIndexFilename)))
	if err := indexTemplate.ExecuteTemplate(w, frontendIndexFilename, page{
		CsrfToken:     csrf.Token(r),
		Title:         getPageTitle(r),
		Description:   getDescription(r),
		OpenGraphType: getOpenGraphType(r),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serve404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	serveIndexPage(w, r)
}

// getPageTitle returns the <title> value of the page. By default it's
// "What Got Done" but if the date or username are available, we prepend those
// to the title, so that it can be "username's What Got Done for the week of
// MMM. D, YYYY".
func getPageTitle(r *http.Request) string {
	t := "What Got Done"

	username, err := usernameFromRequestPath(r)
	if err != nil {
		return t
	}

	date, err := dateFromRequestPath(r)
	if err != nil {
		return t
	}

	dateParsed, err := time.Parse("2006-01-02", string(date))
	if err != nil {
		return t
	}

	formattedDate := dateParsed.Format("Jan. 2, 2006")

	return fmt.Sprintf("%s's What Got Done for the week of %s", username, formattedDate)
}

func getOpenGraphType(r *http.Request) string {
	t := "website"

	_, err := usernameFromRequestPath(r)
	if err != nil {
		return t
	}

	_, err = dateFromRequestPath(r)
	if err != nil {
		return t
	}

	return "article"
}

func getDescription(r *http.Request) string {
	t := "The simple, easy way to share progress with your teammates."

	username, err := usernameFromRequestPath(r)
	if err != nil {
		return t
	}

	date, err := dateFromRequestPath(r)
	if err != nil {
		return t
	}

	dateParsed, err := time.Parse("2006-01-02", string(date))
	if err != nil {
		return t
	}

	formattedDate := dateParsed.Format("January 2, 2006")

	return fmt.Sprintf("Find out what %s accomplished for the week ending on %s", username, formattedDate)
}
