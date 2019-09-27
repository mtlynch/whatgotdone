package handlers

import (
	"net/http"
	"path"
	"time"
)

const (
	userKitAuthCookieName = "userkit_auth_token"
	frontendRootDir = "./frontend/dist"
	frontendIndexFilename = "index.html"
)

var (
	frontendIndexPath = path.Join(frontendRootDir, frontendIndexFilename)
)

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
