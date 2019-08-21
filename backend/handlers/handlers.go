package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const userKitAuthCookieName = "userkit_auth_token"

type recentEntry struct {
	Author       string `json:"author"`
	Date         string `json:"date"`
	lastModified string
	Markdown     string `json:"markdown"`
}

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

func validateEntryDate(date string) bool {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	const whatGotDoneEpochYear = 2019
	if t.Year() < whatGotDoneEpochYear {
		return false
	}
	if t.Weekday() != time.Friday {
		return false
	}
	if t.After(thisFriday()) {
		return false
	}
	return true
}

func thisFriday() time.Time {
	t := time.Now()
	for t.Weekday() != time.Friday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}

func (s defaultServer) loggedInUser(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie(userKitAuthCookieName)
	if err != nil {
		return "", err
	}
	return s.authenticator.UserFromAuthToken(tokenCookie.Value)
}

func usernameFromRequestPath(r *http.Request) (string, error) {
	username := mux.Vars(r)["username"]
	// If something goes wrong in a JavaScript-based client, it will send the literal string "undefined" as the username
	// when the username variable is undefined.
	if username == "undefined" || username == "" {
		return "", errors.New("Invalid username")
	}
	return username, nil
}

func dateFromRequestPath(r *http.Request) (string, error) {
	date := mux.Vars(r)["date"]
	dateFormat := "2006-01-02"
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", errors.New("Invalid date format: must be YYYY-MM-DD")
	}
	return date, nil
}
