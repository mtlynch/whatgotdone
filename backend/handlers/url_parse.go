package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/handlers/parse"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func usernameFromRequestPath(r *http.Request) (types.Username, error) {
	return parse.Username(mux.Vars(r)["username"])
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
