package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/gorilla/csrf"
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
			// If there's no static file that matches this route, serve the index
			// page and let the frontend handle it.
			serveIndexPage(w, r)
			return
		} else if err != nil {
			log.Printf("Failed to retrieve the file %s from the file system: %s", r.URL.Path, err)
			http.Error(w, "Failed to find file: "+r.URL.Path, http.StatusInternalServerError)
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
			// If the client requested a directory, serve the index page.
			serveIndexPage(w, r)
			return
		}

		// Otherwise, serve a static file.
		http.ServeFile(w, r, path.Join(frontendRootDir, r.URL.Path))
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

// getPageTitle returns the <title> value of the page. By default it's
// "What Got Done" but if the date or username are available, we prepend those
// to the title, so that it can be "username's What Got Done for the week of YYYY-MM-DD".
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

	return fmt.Sprintf("%s's What Got Done for the week of %s", username, date)
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
	return fmt.Sprintf("Find out what %s accomplished for the week of %s", username, date)
}
