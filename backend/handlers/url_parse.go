package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/handlers/validate"
)

func usernameFromRequestPath(r *http.Request) (string, error) {
	username := mux.Vars(r)["username"]
	if !validate.Username(username) {
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

func topicFromRequestPath(r *http.Request) (string, error) {
	return mux.Vars(r)["topic"], nil
}
