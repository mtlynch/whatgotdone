package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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
