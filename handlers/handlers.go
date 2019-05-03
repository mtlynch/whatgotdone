package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	userkit "github.com/workpail/userkit-go"

	"github.com/mtlynch/whatgotdone/datastore"
	"github.com/mtlynch/whatgotdone/types"
)

func (s *defaultServer) indexHandler() http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./web/frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		enableCsp(&w)

		type page struct {
			Title string
		}
		p := &page{
			Title: "What Got Done",
		}
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s defaultServer) submitPageHandler() http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./web/frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		enableCsp(&w)

		_, err := s.loggedInUser(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		type page struct {
			Title string
		}
		p := &page{
			Title: "What Got Done - Submit Entry",
		}
		err = templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *defaultServer) entriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		entries, err := s.datastore.All(usernameFromRequestPath(r))
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			return
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) entryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		date, err := dateFromRequestPath(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		j, err := s.datastore.Get(usernameFromRequestPath(r), date)
		if err != nil {
			if _, ok := err.(datastore.EntryNotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Printf("Failed to retrieve entry: %s", err)
			http.Error(w, "Failed to retrieve entry", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(j); err != nil {
			panic(err)
		}
	}
}

func (s defaultServer) userMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		user, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to retrieve information about your account", http.StatusForbidden)
			return
		}

		type userMeResponse struct {
			Username string `json:"username"`
		}

		resp := userMeResponse{
			Username: user.Username,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) submitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == "OPTIONS" {
			return
		}

		user, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to submit a journal entry", http.StatusForbidden)
			return
		}

		decoder := json.NewDecoder(r.Body)

		type submitRequest struct {
			Date         string `json:"date"`
			EntryContent string `json:"entryContent"`
		}

		type submitResponse struct {
			Ok bool `json:"ok"`
		}

		var t submitRequest
		err = decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", r.Body)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
		}
		j := types.JournalEntry{
			Date:         t.Date,
			LastModified: time.Now().Format(time.RFC3339),
			Markdown:     t.EntryContent,
		}
		err = s.datastore.Insert(user.Username, j)
		if err != nil {
			log.Printf("Failed to insert journal entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
		}
		resp := submitResponse{
			Ok: true,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func (s defaultServer) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "userkit_auth_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})

		w.Write([]byte("You are now logged out"))
	}
}

// Catchall for when no API route matches.
func (s *defaultServer) apiRootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// TODO: Adjust this so it's only the CSP for the /login route.
func enableCsp(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Security-Policy", "default-src 'self' https://widget.userkit.io https://api.userkit.io https://www.google.com/recaptcha/api.js https://www.gstatic.com/recaptcha/api2/ https://fonts.googleapis.com")
}

func (s defaultServer) loggedInUser(r *http.Request) (*userkit.User, error) {
	tokenCookie, err := r.Cookie("userkit_auth_token")
	if err != nil {
		return nil, err
	}
	user, err := s.userKitClient.Users.GetUserBySession(tokenCookie.Value)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func usernameFromRequestPath(r *http.Request) string {
	return mux.Vars(r)["username"]
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
