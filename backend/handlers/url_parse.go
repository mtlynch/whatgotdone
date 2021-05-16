package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/handlers/validate"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func usernameFromRequestPath(r *http.Request) (types.Username, error) {
	username := mux.Vars(r)["username"]
	if !validate.Username(username) {
		return "", errors.New("Invalid username")
	}
	return types.Username(username), nil
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

func projectFromRequestPath(r *http.Request) (string, error) {
	return mux.Vars(r)["project"], nil
}
