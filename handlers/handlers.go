package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"

	"github.com/mtlynch/whatgotdone/datastore"
	"github.com/mtlynch/whatgotdone/types"
)

const userKitAuthCookieName = "userkit_auth_token"

type recentEntry struct {
	Author       string `json:"author"`
	Date         string `json:"date"`
	lastModified string
	Markdown     string `json:"markdown"`
}

func (s *defaultServer) indexHandler() http.HandlerFunc {
	var templates = template.Must(
		// Use custom delimiters so Go's delimiters don't clash with Vue's.
		template.New("index.html").Delims("[[", "]]").ParseFiles(
			"./web/frontend/dist/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
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
		username, err := usernameFromRequestPath(r)
		if err != nil {
			log.Printf("Failed to retrieve username from request path: %s", err)
			http.Error(w, "Invalid username", http.StatusBadRequest)
		}

		entries, err := s.datastore.All(username)
		if err != nil {
			log.Printf("Failed to retrieve entries: %s", err)
			http.Error(w, fmt.Sprintf("Failed to retrieve entries for %s", username), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) recentEntriesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.datastore.Users()
		if err != nil {
			log.Printf("Failed to retrieve users: %s", err)
			return
		}

		var entries []recentEntry
		for _, username := range users {
			userEntries, err := s.datastore.All(username)
			if err != nil {
				log.Printf("Failed to retrieve entries for user %s: %s", username, err)
				return
			}
			for _, entry := range userEntries {
				entries = append(entries, recentEntry{
					Author:       username,
					Date:         entry.Date,
					lastModified: entry.LastModified,
					Markdown:     entry.Markdown,
				})
			}
		}

		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Date < entries[j].Date {
				return true
			}
			if entries[i].Date > entries[j].Date {
				return false
			}
			return entries[i].lastModified < entries[j].lastModified
		})

		// Reverse the order of entries.
		for i := len(entries)/2 - 1; i >= 0; i-- {
			opp := len(entries) - 1 - i
			entries[i], entries[opp] = entries[opp], entries[i]
		}

		const maxEntries = 10
		if len(entries) > maxEntries {
			entries = entries[:maxEntries]
		}

		if err := json.NewEncoder(w).Encode(entries); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) draftHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
		} else if r.Method == "GET" {
			s.handleDraftGet(w, r)
		} else if r.Method == "POST" {
			s.handleDraftPost(w, r)
		} else {
			log.Printf("Invalid method for drafts handler: %s", r.Method)
			http.Error(w, "Invalid operation", http.StatusBadRequest)
		}
	}
}

func (s defaultServer) handleDraftGet(w http.ResponseWriter, r *http.Request) {
	username, err := s.loggedInUser(r)
	if err != nil {
		http.Error(w, "You must log in to retrieve a draft entry", http.StatusForbidden)
		return
	}

	date, err := dateFromRequestPath(r)
	if err != nil {
		log.Printf("Invalid date: %s", date)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	j, err := s.datastore.GetDraft(username, date)
	if _, ok := err.(datastore.DraftNotFoundError); ok {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Failed to retrieve draft entry: %s", err)
		http.Error(w, "Failed to retrieve draft entry", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(j); err != nil {
		panic(err)
	}
}

func (s defaultServer) handleDraftPost(w http.ResponseWriter, r *http.Request) {
	username, err := s.loggedInUser(r)
	if err != nil {
		http.Error(w, "You must log in to save a draft entry", http.StatusForbidden)
		return
	}

	type draftRequest struct {
		EntryContent string `json:"entryContent"`
	}

	type draftResponse struct {
		Ok bool `json:"ok"`
	}

	var t draftRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&t)
	if err != nil {
		log.Printf("Failed to decode request: %s", err)
		http.Error(w, "Failed to decode request", http.StatusBadRequest)
	}

	date, err := dateFromRequestPath(r)
	if err != nil {
		log.Printf("Invalid date: %s", date)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j := types.JournalEntry{
		Date:         date,
		LastModified: time.Now().Format(time.RFC3339),
		Markdown:     t.EntryContent,
	}
	err = s.datastore.InsertDraft(username, j)
	if err != nil {
		log.Printf("Failed to update draft entry: %s", err)
		http.Error(w, "Failed to update draft entry", http.StatusInternalServerError)
		return
	}
	resp := draftResponse{
		Ok: true,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		panic(err)
	}
}

func (s defaultServer) userMeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must be logged in to retrieve information about your account", http.StatusForbidden)
			return
		}

		type userMeResponse struct {
			Username string `json:"username"`
		}

		resp := userMeResponse{
			Username: username,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	}
}

func (s *defaultServer) submitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			return
		}

		username, err := s.loggedInUser(r)
		if err != nil {
			http.Error(w, "You must log in to submit a journal entry", http.StatusForbidden)
			return
		}

		type submitRequest struct {
			Date         string `json:"date"`
			EntryContent string `json:"entryContent"`
		}

		type submitResponse struct {
			Ok   bool   `json:"ok"`
			Path string `json:"path"`
		}

		var t submitRequest
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&t)
		if err != nil {
			log.Printf("Failed to decode request: %s", err)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
		}
		if !validateEntryDate(t.Date) {
			log.Printf("Invalid date: %s", t.Date)
			http.Error(w, "Invalid date", http.StatusBadRequest)
			return
		}

		j := types.JournalEntry{
			Date:         t.Date,
			LastModified: time.Now().Format(time.RFC3339),
			Markdown:     t.EntryContent,
		}

		err = s.datastore.InsertDraft(username, j)
		if err != nil {
			log.Printf("Failed to update journal draft entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
			return
		}
		err = s.datastore.Insert(username, j)
		if err != nil {
			log.Printf("Failed to insert journal entry: %s", err)
			http.Error(w, "Failed to insert entry", http.StatusInternalServerError)
			return
		}
		resp := submitResponse{
			Ok:   true,
			Path: fmt.Sprintf("/%s/%s", username, t.Date),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
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
		http.Error(w, "Invalid API path", http.StatusBadRequest)
	}
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
