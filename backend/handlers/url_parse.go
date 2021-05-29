package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mtlynch/whatgotdone/backend/handlers/parse"
	"github.com/mtlynch/whatgotdone/backend/types"
)

func usernameFromRequestPath(r *http.Request) (types.Username, error) {
	return parse.Username(mux.Vars(r)["username"])
}

func dateFromRequestPath(r *http.Request) (types.EntryDate, error) {
	return parse.EntryDate(mux.Vars(r)["date"])
}

func projectFromRequestPath(r *http.Request) (string, error) {
	return mux.Vars(r)["project"], nil
}
